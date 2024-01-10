// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"context"
	"fmt"
	"time"

	"github.com/ledgerwatch/erigon-lib/chain"
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/log/v3"

	"github.com/ledgerwatch/erigon/consensus"
	"github.com/ledgerwatch/erigon/core/state"
	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/core/vm"
	"github.com/ledgerwatch/erigon/core/vm/evmtypes"
	"github.com/ledgerwatch/erigon/crypto"
	"github.com/ledgerwatch/erigon/rpc"
)

// applyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func applyTransaction(config *chain.Config, engine consensus.EngineReader, gp *GasPool, ibs *state.IntraBlockState,
	stateWriter state.StateWriter, header *types.Header, tx types.Transaction, usedGas, usedBlobGas *uint64,
	evm *vm.EVM, cfg vm.Config, historicalRPCService *rpc.Client, historicalRPCTimeout *time.Duration) (*types.Receipt, []byte, error) {
	rules := evm.ChainRules()
	msg, err := tx.AsMessage(*types.MakeSigner(config, header.Number.Uint64(), header.Time), header.BaseFee, rules)
	if err != nil {
		return nil, nil, err
	}
	msg.SetCheckNonce(!cfg.StatelessExec)

	if msg.FeeCap().IsZero() && engine != nil {
		// Only zero-gas transactions may be service ones
		syscall := func(contract libcommon.Address, data []byte) ([]byte, error) {
			return SysCallContract(contract, data, config, ibs, header, engine, true /* constCall */)
		}
		msg.SetIsFree(engine.IsServiceTransaction(msg.From(), syscall))
	}

	txContext := NewEVMTxContext(msg)
	if cfg.TraceJumpDest {
		txContext.TxHash = tx.Hash()
	}

	// Update the evm with the new transaction context.
	evm.Reset(txContext, ibs)

	nonce := tx.GetNonce()
	if msg.IsDepositTx() && config.IsOptimismRegolith(evm.Context().Time) {
		nonce = ibs.GetNonce(msg.From())
	}

	var legacyReceipt *types.Receipt
	isOptimismPreBlock := evm.ChainConfig().IsOptimismPreBedrock(header.Number.Uint64())
	result := &ExecutionResult{
		UsedGas:    msg.Gas(),
		Err:        nil,
		ReturnData: []byte{},
	}

	if !isOptimismPreBlock {
		result, err = ApplyMessage(evm, msg, gp, true /* refunds */, false /* gasBailout */)
		if err != nil {
			return nil, nil, err
		}
	}

	// Update the state with pending changes
	if err = ibs.FinalizeTx(rules, stateWriter); err != nil {
		return nil, nil, err
	}
	*usedGas += result.UsedGas
	if usedBlobGas != nil {
		*usedBlobGas += tx.GetBlobGas()
	}

	if isOptimismPreBlock {
		if historicalRPCService != nil {
			ctx, cancel := context.WithTimeout(context.Background(), *historicalRPCTimeout)
			err = historicalRPCService.CallContext(ctx, &legacyReceipt, "eth_getTransactionReceipt", tx.Hash().String())
			cancel()
			if err != nil {
				return nil, nil, err
			}
			*usedGas = legacyReceipt.GasUsed
		} else {
			return nil, nil, fmt.Errorf("legacy block must be handled by the historicalRPCService")
		}
	}

	// Set the receipt logs and create the bloom filter.
	// based on the eip phase, we're passing whether the root touch-delete accounts.
	var receipt *types.Receipt
	if !cfg.NoReceipts {
		// by the tx.
		receipt = &types.Receipt{Type: tx.Type(), CumulativeGasUsed: *usedGas}
		if result.Failed() {
			receipt.Status = types.ReceiptStatusFailed
		} else {
			receipt.Status = types.ReceiptStatusSuccessful
		}
		receipt.TxHash = tx.Hash()
		receipt.GasUsed = result.UsedGas

		if msg.IsDepositTx() && config.IsOptimismRegolith(evm.Context().Time) {
			// The actual nonce for deposit transactions is only recorded from Regolith onwards and
			// otherwise must be nil.
			receipt.DepositNonce = &nonce

			if config.IsOptimismCanyon(evm.Context().Time) {
				receipt.DepositReceiptVersion = new(uint64)
				*receipt.DepositReceiptVersion = types.CanyonDepositReceiptVersion
			}
		}

		// if the transaction created a contract, store the creation address in the receipt.
		if msg.To() == nil {
			receipt.ContractAddress = crypto.CreateAddress(evm.TxContext().Origin, nonce)
		}
		// Set the receipt logs and create a bloom for filtering
		receipt.Logs = ibs.GetLogs(tx.Hash())
		receipt.BlockNumber = header.Number
		receipt.TransactionIndex = uint(ibs.TxIndex())

		if isOptimismPreBlock {
			receipt.GasUsed = legacyReceipt.GasUsed
			receipt.Logs = legacyReceipt.Logs
			receipt.Status = legacyReceipt.Status
			// The following fields can not be set in the legacy receipt
			// The math of calculating legacy and new receipts is not compatible
			receipt.L1GasPrice = legacyReceipt.L1GasPrice
			receipt.L1GasUsed = legacyReceipt.L1GasUsed
			receipt.L1Fee = legacyReceipt.L1Fee
			receipt.FeeScalar = legacyReceipt.FeeScalar
		}

		receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	}

	return receipt, result.ReturnData, err
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *chain.Config, blockHashFunc func(n uint64) libcommon.Hash, engine consensus.EngineReader,
	author *libcommon.Address, gp *GasPool, ibs *state.IntraBlockState, stateWriter state.StateWriter,
	header *types.Header, tx types.Transaction, usedGas, usedBlobGas *uint64, cfg vm.Config,
	historicalRPCService *rpc.Client, historicalRPCTimeout *time.Duration,
) (*types.Receipt, []byte, error) {
	log.Debug("ApplyTransaction called for", "txhash", tx.Hash(), "blockNum", header.Number.Uint64())
	// Create a new context to be used in the EVM environment

	// Add addresses to access list if applicable
	// about the transaction and calling mechanisms.
	cfg.SkipAnalysis = SkipAnalysis(config, header.Number.Uint64())

	l1CostFunc := types.NewL1CostFunc(config, ibs)
	blockContext := NewEVMBlockContext(header, blockHashFunc, engine, author, l1CostFunc)
	vmenv := vm.NewEVM(blockContext, evmtypes.TxContext{}, ibs, config, cfg)

	return applyTransaction(config, engine, gp, ibs, stateWriter, header, tx, usedGas, usedBlobGas, vmenv, cfg, historicalRPCService, historicalRPCTimeout)
}
