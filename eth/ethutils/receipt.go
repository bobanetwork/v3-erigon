package ethutils

import (
	"math/big"

	"github.com/erigontech/erigon-lib/log/v3"
	"github.com/holiman/uint256"

	"github.com/erigontech/erigon-lib/chain"
	"github.com/erigontech/erigon-lib/common"
	"github.com/erigontech/erigon-lib/common/hexutil"
	"github.com/erigontech/erigon/consensus/misc"
	"github.com/erigontech/erigon/core/types"
)

func MarshalReceipt(
	receipt *types.Receipt,
	txn types.Transaction,
	chainConfig *chain.Config,
	header *types.Header,
	txnHash common.Hash,
	signed bool,
) map[string]interface{} {
	var chainId *big.Int
	switch t := txn.(type) {
	case *types.LegacyTx:
		if t.Protected() {
			chainId = types.DeriveChainId(&t.V).ToBig()
		}
	case *types.DepositTx:
		// Deposit TX does not have chain ID
	default:
		chainId = txn.GetChainID().ToBig()
	}

	var from common.Address
	if signed {
		signer := types.LatestSignerForChainID(chainId)
		from, _ = txn.Sender(*signer)
	}

	fields := map[string]interface{}{
		"blockHash":         receipt.BlockHash,
		"blockNumber":       hexutil.Uint64(receipt.BlockNumber.Uint64()),
		"transactionHash":   txnHash,
		"transactionIndex":  hexutil.Uint64(receipt.TransactionIndex),
		"from":              from,
		"to":                txn.GetTo(),
		"type":              hexutil.Uint(txn.Type()),
		"gasUsed":           hexutil.Uint64(receipt.GasUsed),
		"cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed),
		"contractAddress":   nil,
		"logs":              receipt.Logs,
		"logsBloom":         types.CreateBloom(types.Receipts{receipt}),
	}

	if !chainConfig.IsLondon(header.Number.Uint64()) {
		fields["effectiveGasPrice"] = (*hexutil.Big)(txn.GetPrice().ToBig())
	} else {
		baseFee, _ := uint256.FromBig(header.BaseFee)
		gasPrice := new(big.Int).Add(header.BaseFee, txn.GetEffectiveGasTip(baseFee).ToBig())
		fields["effectiveGasPrice"] = (*hexutil.Big)(gasPrice)
	}

	// Assign receipt status.
	fields["status"] = hexutil.Uint64(receipt.Status)
	if receipt.Logs == nil {
		fields["logs"] = [][]*types.Log{}
	}

	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if receipt.ContractAddress != (common.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	}

	if chainConfig.IsOptimism() {
		if txn.Type() != types.DepositTxType {
			fields["l1GasPrice"] = hexutil.Big(*receipt.L1GasPrice)
			fields["l1GasUsed"] = hexutil.Big(*receipt.L1GasUsed)
			fields["l1Fee"] = hexutil.Big(*receipt.L1Fee)
			if receipt.FeeScalar != nil { // removed in Ecotone
				fields["l1FeeScalar"] = receipt.FeeScalar
			}
			if receipt.L1BaseFeeScalar != nil { // added in Ecotone
				fields["l1BaseFeeScalar"] = hexutil.Uint64(*receipt.L1BaseFeeScalar)
			}
			if receipt.L1BlobBaseFee != nil { // added in Ecotone
				fields["l1BlobBaseFee"] = hexutil.Big(*receipt.L1BlobBaseFee)
			}
			if receipt.L1BlobBaseFeeScalar != nil { // added in Ecotone
				fields["l1BlobBaseFeeScalar"] = hexutil.Uint64(*receipt.L1BlobBaseFeeScalar)
			}
		} else {
			if receipt.DepositNonce != nil {
				fields["depositNonce"] = hexutil.Uint64(*receipt.DepositNonce)
			}
			if receipt.DepositReceiptVersion != nil {
				fields["depositReceiptVersion"] = hexutil.Uint64(*receipt.DepositReceiptVersion)
			}
		}
	}

	// Set derived blob related fields
	numBlobs := len(txn.GetBlobHashes())
	if numBlobs > 0 {
		if header.ExcessBlobGas == nil {
			log.Warn("excess blob gas not set when trying to marshal blob tx")
		} else {
			blobGasPrice, err := misc.GetBlobGasPrice(chainConfig, *header.ExcessBlobGas)
			if err != nil {
				log.Error(err.Error())
			}
			fields["blobGasPrice"] = (*hexutil.Big)(blobGasPrice.ToBig())
			fields["blobGasUsed"] = hexutil.Uint64(misc.GetBlobGasUsed(numBlobs))
		}
	}

	return fields
}
