package state

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/erigontech/erigon-lib/log/v3"
	"github.com/stretchr/testify/require"

	"github.com/erigontech/erigon-lib/common"
	"github.com/erigontech/erigon-lib/recsplit"
	"github.com/erigontech/erigon-lib/seg"
)

func Benchmark_BtreeIndex_Allocation(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < b.N; i++ {
		now := time.Now()
		count := rnd.Intn(1000000000)
		bt := newBtAlloc(uint64(count), uint64(1<<12), true)
		bt.traverseDfs()
		fmt.Printf("alloc %v\n", time.Since(now))
	}
}

func Benchmark_BtreeIndex_Search(b *testing.B) {
	logger := log.New()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tmp := b.TempDir()
	defer os.RemoveAll(tmp)
	dataPath := "../../data/storage.256-288.kv"

	indexPath := path.Join(tmp, filepath.Base(dataPath)+".bti")
	err := BuildBtreeIndex(dataPath, indexPath, logger)
	require.NoError(b, err)

	M := 1024
	bt, err := OpenBtreeIndex(indexPath, dataPath, uint64(M))

	require.NoError(b, err)

	idx := NewBtIndexReader(bt)

	keys, err := pivotKeysFromKV(dataPath)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		p := rnd.Intn(len(keys))
		cur, err := idx.Seek(keys[p])
		require.NoErrorf(b, err, "i=%d", i)
		require.EqualValues(b, keys[p], cur.key)
		require.NotEmptyf(b, cur.Value(), "i=%d", i)
	}

	bt.Close()
}

func benchInitBtreeIndex(b *testing.B, M uint64) (*BtIndex, [][]byte, string) {
	b.Helper()

	logger := log.New()
	tmp := b.TempDir()
	b.Cleanup(func() { os.RemoveAll(tmp) })

	dataPath := generateCompressedKV(b, tmp, 52, 10, 1000000, logger)
	indexPath := path.Join(tmp, filepath.Base(dataPath)+".bt")
	bt, err := CreateBtreeIndex(indexPath, dataPath, M, logger)
	require.NoError(b, err)

	keys, err := pivotKeysFromKV(dataPath)
	require.NoError(b, err)
	return bt, keys, dataPath
}

func Benchmark_BTree_Seek(b *testing.B) {
	M := uint64(1024)
	bt, keys, _ := benchInitBtreeIndex(b, M)
	defer bt.Close()

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.Run("seek_only", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p := rnd.Intn(len(keys))

			cur, err := bt.Seek(keys[p])
			require.NoError(b, err)

			require.EqualValues(b, keys[p], cur.key)
		}
	})

	b.Run("seek_then_next", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p := rnd.Intn(len(keys))

			cur, err := bt.Seek(keys[p])
			require.NoError(b, err)

			require.EqualValues(b, keys[p], cur.key)

			prevKey := common.Copy(keys[p])
			ntimer := time.Duration(0)
			nextKeys := 5000
			for j := 0; j < nextKeys; j++ {
				ntime := time.Now()

				if !cur.Next() {
					break
				}
				ntimer += time.Since(ntime)

				nk := cur.Key()
				if bytes.Compare(prevKey, nk) > 0 {
					b.Fatalf("prev %s cur %s, next key should be greater", prevKey, nk)
				}
				prevKey = nk
			}
			if i%1000 == 0 {
				fmt.Printf("next_access_last[of %d keys] %v\n", nextKeys, ntimer/time.Duration(nextKeys))
			}

		}
	})
}

// requires existing KV index file at ../../data/storage.kv
func Benchmark_Recsplit_Find_ExternalFile(b *testing.B) {
	dataPath := "../../data/storage.kv"
	f, err := os.Stat(dataPath)
	if err != nil || f.IsDir() {
		b.Skip("requires existing KV index file at ../../data/storage.kv")
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tmp := b.TempDir()

	defer os.RemoveAll(tmp)

	indexPath := dataPath + "i"
	idx, err := recsplit.OpenIndex(indexPath)
	require.NoError(b, err)
	idxr := recsplit.NewIndexReader(idx)

	decomp, err := seg.NewDecompressor(dataPath)
	require.NoError(b, err)
	defer decomp.Close()

	getter := decomp.MakeGetter()

	keys, err := pivotKeysFromKV(dataPath)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		p := rnd.Intn(len(keys))

		offset, _ := idxr.Lookup(keys[p])
		getter.Reset(offset)

		require.True(b, getter.HasNext())

		key, pa := getter.Next(nil)
		require.NotEmpty(b, key)

		value, pb := getter.Next(nil)
		if pb-pa != 1 {
			require.NotEmpty(b, value)
		}

		require.NoErrorf(b, err, "i=%d", i)
		require.EqualValues(b, keys[p], key)
	}
}
