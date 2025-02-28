package jsonrpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/erigontech/erigon-lib/common"
	"github.com/erigontech/erigon-lib/common/hexutil"
	"github.com/erigontech/erigon-lib/kv"
	"github.com/erigontech/erigon/core/rawdb"
	"github.com/erigontech/erigon/core/types"
	"github.com/erigontech/erigon/eth/ethutils"
	"github.com/erigontech/erigon/rpc"
	"github.com/erigontech/erigon/turbo/adapter/ethapi"
	"github.com/erigontech/erigon/turbo/rpchelper"
)

type GraphQLAPI interface {
	GetBlockDetails(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error)
	GetChainID(ctx context.Context) (*big.Int, error)
}

type GraphQLAPIImpl struct {
	*BaseAPI
	db kv.RoDB
}

func NewGraphQLAPI(base *BaseAPI, db kv.RoDB) *GraphQLAPIImpl {
	return &GraphQLAPIImpl{
		BaseAPI: base,
		db:      db,
	}
}

func (api *GraphQLAPIImpl) GetChainID(ctx context.Context) (*big.Int, error) {
	tx, err := api.db.BeginRo(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	response, err := api.chainConfig(ctx, tx)
	if err != nil {
		return nil, err
	}

	return response.ChainID, nil
}

func (api *GraphQLAPIImpl) GetBlockDetails(ctx context.Context, blockNumber rpc.BlockNumber) (map[string]interface{}, error) {
	tx, err := api.db.BeginRo(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	block, senders, err := api.getBlockWithSenders(ctx, blockNumber, tx)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, nil
	}

	getBlockRes, err := api.delegateGetBlockByNumber(tx, block, blockNumber, false)
	if err != nil {
		return nil, err
	}

	chainConfig, err := api.chainConfig(ctx, tx)
	if err != nil {
		return nil, err
	}

	receipts, err := api.getReceipts(ctx, tx, block, senders)
	if err != nil {
		return nil, fmt.Errorf("getReceipts error: %w", err)
	}

	result := make([]map[string]interface{}, 0, len(receipts))
	for _, receipt := range receipts {
		txn := block.Transactions()[receipt.TransactionIndex]

		transaction := ethutils.MarshalReceipt(receipt, txn, chainConfig, block.HeaderNoCopy(), txn.Hash(), true)
		transaction["nonce"] = txn.GetNonce()
		transaction["value"] = txn.GetValue()
		transaction["data"] = txn.GetData()
		transaction["logs"] = receipt.Logs
		result = append(result, transaction)
	}

	response := map[string]interface{}{}
	response["block"] = getBlockRes
	response["receipts"] = result

	return response, nil
}

func (api *GraphQLAPIImpl) getBlockWithSenders(ctx context.Context, number rpc.BlockNumber, tx kv.Tx) (*types.Block, []common.Address, error) {
	if number == rpc.PendingBlockNumber {
		return api.pendingBlock(), nil, nil
	}

	blockHeight, blockHash, _, err := rpchelper.GetBlockNumber(rpc.BlockNumberOrHashWithNumber(number), tx, api.filters)
	if err != nil {
		return nil, nil, err
	}

	block, err := api.blockWithSenders(ctx, tx, blockHash, blockHeight)
	if err != nil {
		return nil, nil, err
	}
	if block == nil {
		return nil, nil, nil
	}
	return block, block.Body().SendersFromTxs(), nil
}

func (api *GraphQLAPIImpl) delegateGetBlockByNumber(tx kv.Tx, b *types.Block, number rpc.BlockNumber, inclTx bool) (map[string]interface{}, error) {
	td, err := rawdb.ReadTd(tx, b.Hash(), b.NumberU64())
	if err != nil {
		return nil, err
	}
	additionalFields := make(map[string]interface{})
	receipts := rawdb.ReadRawReceipts(tx, uint64(number.Int64()))
	response, err := ethapi.RPCMarshalBlock(b, inclTx, inclTx, additionalFields, receipts)
	if !inclTx {
		delete(response, "transactions") // workaround for https://github.com/erigontech/erigon/issues/4989#issuecomment-1218415666
	}
	response["totalDifficulty"] = (*hexutil.Big)(td)
	response["transactionCount"] = b.Transactions().Len()

	if err == nil && number == rpc.PendingBlockNumber {
		// Pending blocks need to nil out a few fields
		for _, field := range []string{"hash", "nonce", "miner"} {
			response[field] = nil
		}
	}

	return response, err
}
