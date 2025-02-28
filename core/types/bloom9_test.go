// Copyright 2014 The go-ethereum Authors
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

package types

import (
	"fmt"
	"math/big"
	"testing"

	libcommon "github.com/erigontech/erigon-lib/common"
	"github.com/holiman/uint256"

	"github.com/erigontech/erigon/crypto"
)

func TestBloom(t *testing.T) {
	t.Parallel()
	positive := []string{
		"testtest",
		"test",
		"hallo",
		"other",
	}
	negative := []string{
		"tes",
		"lo",
	}

	var bloom Bloom
	for _, data := range positive {
		bloom.Add([]byte(data))
	}

	for _, data := range positive {
		if !bloom.Test([]byte(data)) {
			t.Error("expected", data, "to test true")
		}
	}
	for _, data := range negative {
		if bloom.Test([]byte(data)) {
			t.Error("did not expect", data, "to test true")
		}
	}
}

// TestBloomExtensively does some more thorough tests
func TestBloomExtensively(t *testing.T) {
	t.Parallel()
	var exp = libcommon.HexToHash("c8d3ca65cdb4874300a9e39475508f23ed6da09fdbc487f89a2dcf50b09eb263")
	var b Bloom
	// Add 100 "random" things
	for i := 0; i < 100; i++ {
		data := fmt.Sprintf("xxxxxxxxxx data %d yyyyyyyyyyyyyy", i)
		b.Add([]byte(data))
		//b.Add(new(big.Int).SetBytes([]byte(data)))
	}
	got := crypto.Keccak256Hash(b.Bytes())
	if got != exp {
		t.Errorf("Got %x, exp %x", got, exp)
	}
	var b2 Bloom
	b2.SetBytes(b.Bytes())
	got2 := crypto.Keccak256Hash(b2.Bytes())
	if got != got2 {
		t.Errorf("Got %x, exp %x", got, got2)
	}
}

func BenchmarkBloom9(b *testing.B) {
	test := []byte("testestestest")
	for i := 0; i < b.N; i++ {
		Bloom9(test)
	}
}

func BenchmarkBloom9Lookup(b *testing.B) {
	toTest := []byte("testtest")
	bloom := new(Bloom)
	for i := 0; i < b.N; i++ {
		bloom.Test(toTest)
	}
}

func BenchmarkCreateBloom(b *testing.B) {

	one, _ := uint256.FromBig(big.NewInt(1))
	two, _ := uint256.FromBig(big.NewInt(2))

	var txs = Transactions{
		NewContractCreation(1, one, 1, one, nil),
		NewTransaction(2, libcommon.HexToAddress("0x2"), two, 2, two, nil),
	}
	var rSmall = Receipts{
		&Receipt{
			Status:            ReceiptStatusFailed,
			CumulativeGasUsed: 1,
			Logs: []*Log{
				{Address: libcommon.BytesToAddress([]byte{0x11})},
				{Address: libcommon.BytesToAddress([]byte{0x01, 0x11})},
			},
			TxHash:          txs[0].Hash(),
			ContractAddress: libcommon.BytesToAddress([]byte{0x01, 0x11, 0x11}),
			GasUsed:         1,
		},
		&Receipt{
			PostState:         libcommon.Hash{2}.Bytes(),
			CumulativeGasUsed: 3,
			Logs: []*Log{
				{Address: libcommon.BytesToAddress([]byte{0x22})},
				{Address: libcommon.BytesToAddress([]byte{0x02, 0x22})},
			},
			TxHash:          txs[1].Hash(),
			ContractAddress: libcommon.BytesToAddress([]byte{0x02, 0x22, 0x22}),
			GasUsed:         2,
		},
	}

	var rLarge = make(Receipts, 200)
	// Fill it with 200 receipts x 2 logs
	for i := 0; i < 200; i += 2 {
		copy(rLarge[i:], rSmall)
	}
	b.Run("small", func(b *testing.B) {
		b.ReportAllocs()
		var bl Bloom
		for i := 0; i < b.N; i++ {
			bl = CreateBloom(rSmall)
		}
		b.StopTimer()
		var exp = libcommon.HexToHash("c384c56ece49458a427c67b90fefe979ebf7104795be65dc398b280f24104949")
		got := crypto.Keccak256Hash(bl.Bytes())
		if got != exp {
			b.Errorf("Got %x, exp %x", got, exp)
		}
	})
	b.Run("large", func(b *testing.B) {
		b.ReportAllocs()
		var bl Bloom
		for i := 0; i < b.N; i++ {
			bl = CreateBloom(rLarge)
		}
		b.StopTimer()
		var exp = libcommon.HexToHash("c384c56ece49458a427c67b90fefe979ebf7104795be65dc398b280f24104949")
		got := crypto.Keccak256Hash(bl.Bytes())
		if got != exp {
			b.Errorf("Got %x, exp %x", got, exp)
		}
	})
}
