package finality

import (
	"github.com/erigontech/erigon-lib/common"
	"github.com/erigontech/erigon-lib/kv"
	"github.com/erigontech/erigon/core/rawdb"
	"github.com/erigontech/erigon/core/types"
	"github.com/erigontech/erigon/polygon/bor/finality/whitelist"
)

func GetFinalizedBlockNumber(tx kv.Tx) uint64 {
	currentBlockNum := rawdb.ReadCurrentHeader(tx)

	service := whitelist.GetWhitelistingService()

	doExist, number, hash := service.GetWhitelistedMilestone()
	if doExist && number <= currentBlockNum.Number.Uint64() {

		blockHeader := rawdb.ReadHeaderByNumber(tx, number)

		if blockHeader == nil {
			return 0
		}

		if blockHeader.Hash() == hash {
			return number
		}
	}

	doExist, number, hash = service.GetWhitelistedCheckpoint()
	if doExist && number <= currentBlockNum.Number.Uint64() {

		blockHeader := rawdb.ReadHeaderByNumber(tx, number)

		if blockHeader == nil {
			return 0
		}

		if blockHeader.Hash() == hash {
			return number
		}
	}

	return 0
}

// CurrentFinalizedBlock retrieves the current finalized block of the canonical
// chain. The block is retrieved from the blockchain's internal cache.
func CurrentFinalizedBlock(tx kv.Tx, number uint64) *types.Block {
	hash, err := rawdb.ReadCanonicalHash(tx, number)
	if err != nil || hash == (common.Hash{}) {
		return nil
	}

	return rawdb.ReadBlock(tx, hash, number)
}
