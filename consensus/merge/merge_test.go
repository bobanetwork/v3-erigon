package merge

import (
	"math/big"
	"testing"

	"github.com/erigontech/erigon-lib/chain"
	libcommon "github.com/erigontech/erigon-lib/common"

	"github.com/erigontech/erigon/consensus"
	"github.com/erigontech/erigon/core/types"
)

type readerMock struct{}

func (r readerMock) Config() *chain.Config {
	return nil
}

func (r readerMock) CurrentHeader() *types.Header {
	return nil
}

func (r readerMock) GetHeader(libcommon.Hash, uint64) *types.Header {
	return nil
}

func (r readerMock) GetHeaderByNumber(uint64) *types.Header {
	return nil
}

func (r readerMock) GetHeaderByHash(libcommon.Hash) *types.Header {
	return nil
}

func (r readerMock) GetTd(libcommon.Hash, uint64) *big.Int {
	return nil
}

func (r readerMock) FrozenBlocks() uint64 {
	return 0
}

func (r readerMock) BorSpan(spanId uint64) []byte {
	return nil
}

// The thing only that changes beetwen normal ethash checks other than POW, is difficulty
// and nonce so we are gonna test those
func TestVerifyHeaderDifficulty(t *testing.T) {
	header := &types.Header{
		Difficulty: big.NewInt(1),
		Time:       1,
	}

	parent := &types.Header{}

	var eth1Engine consensus.Engine
	mergeEngine := New(eth1Engine)

	err := mergeEngine.verifyHeader(readerMock{}, header, parent)
	if err != errInvalidDifficulty {
		if err != nil {
			t.Fatalf("Merge engine should not accept non-zero difficulty, got %s", err.Error())
		} else {
			t.Fatalf("Merge engine should not accept non-zero difficulty")
		}
	}
}

func TestVerifyHeaderNonce(t *testing.T) {
	header := &types.Header{
		Nonce:      types.BlockNonce{1, 0, 0, 0, 0, 0, 0, 0},
		Difficulty: big.NewInt(0),
		Time:       1,
	}

	parent := &types.Header{}

	var eth1Engine consensus.Engine
	mergeEngine := New(eth1Engine)

	err := mergeEngine.verifyHeader(readerMock{}, header, parent)
	if err != errInvalidNonce {
		if err != nil {
			t.Fatalf("Merge engine should not accept non-zero difficulty, got %s", err.Error())
		} else {
			t.Fatalf("Merge engine should not accept non-zero difficulty")
		}
	}
}
