package jsonrpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	libcommon "github.com/erigontech/erigon-lib/common"
	"github.com/erigontech/erigon-lib/direct"
	"github.com/erigontech/erigon-lib/gointerfaces/sentry"
	"github.com/erigontech/erigon-lib/wrap"

	"github.com/erigontech/erigon-lib/log/v3"

	"github.com/erigontech/erigon/cmd/rpcdaemon/rpcservices"
	"github.com/erigontech/erigon/core"
	"github.com/erigontech/erigon/eth/protocols/eth"
	"github.com/erigontech/erigon/ethdb/privateapi"
	"github.com/erigontech/erigon/rlp"
	"github.com/erigontech/erigon/turbo/builder"
	"github.com/erigontech/erigon/turbo/rpchelper"
	"github.com/erigontech/erigon/turbo/stages"
	"github.com/erigontech/erigon/turbo/stages/mock"
)

func TestEthSubscribe(t *testing.T) {
	m, require := mock.Mock(t), require.New(t)
	chain, err := core.GenerateChain(m.ChainConfig, m.Genesis, m.Engine, m.DB, 7, func(i int, b *core.BlockGen) {
		b.SetCoinbase(libcommon.Address{1})
	})
	require.NoError(err)

	b, err := rlp.EncodeToBytes(&eth.BlockHeadersPacket66{
		RequestId:          1,
		BlockHeadersPacket: chain.Headers,
	})
	require.NoError(err)

	m.ReceiveWg.Add(1)
	for _, err = range m.Send(&sentry.InboundMessage{Id: sentry.MessageId_BLOCK_HEADERS_66, Data: b, PeerId: m.PeerId}) {
		require.NoError(err)
	}
	m.ReceiveWg.Wait() // Wait for all messages to be processed before we proceeed

	ctx := context.Background()
	logger := log.New()
	backendServer := privateapi.NewEthBackendServer(ctx, nil, m.DB, m.Notifications.Events, m.BlockReader, logger, builder.NewLatestBlockBuiltStore())
	backendClient := direct.NewEthBackendClientDirect(backendServer)
	backend := rpcservices.NewRemoteBackend(backendClient, m.DB, m.BlockReader)
	ff := rpchelper.New(ctx, rpchelper.DefaultFiltersConfig, backend, nil, nil, func() {}, m.Log)

	newHeads, id := ff.SubscribeNewHeads(16)
	defer ff.UnsubscribeHeads(id)

	initialCycle := mock.MockInsertAsInitialCycle
	highestSeenHeader := chain.TopBlock.NumberU64()

	hook := stages.NewHook(m.Ctx, m.DB, m.Notifications, m.Sync, m.BlockReader, m.ChainConfig, m.Log, nil)
	if err := stages.StageLoopIteration(m.Ctx, m.DB, wrap.TxContainer{}, m.Sync, initialCycle, logger, m.BlockReader, hook, false); err != nil {
		t.Fatal(err)
	}

	for i := uint64(1); i <= highestSeenHeader; i++ {
		header := <-newHeads
		fmt.Printf("Got header %d\n", header.Number.Uint64())
		require.Equal(i, header.Number.Uint64())
	}
}
