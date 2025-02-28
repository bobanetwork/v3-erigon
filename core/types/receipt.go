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
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/erigontech/erigon-lib/chain"
	libcommon "github.com/erigontech/erigon-lib/common"
	"github.com/erigontech/erigon-lib/common/hexutil"
	"github.com/erigontech/erigon-lib/common/hexutility"
	"github.com/erigontech/erigon-lib/opstack"

	"github.com/erigontech/erigon/crypto"
	"github.com/erigontech/erigon/rlp"
)

// go:generate gencodec -type Receipt -field-override receiptMarshaling -out gen_receipt_json.go

// to get it working need to update github.com/ugorji/go/codec to v1.2.12 which has the fix:
//   - https://github.com/ugorji/go/commit/8286c2dc986535d23e3fad8d3e816b9dd1e5aea6
// however updating the lib has caused us issues in the past, and we don't have good unit test coverage for updating atm
// we also use this for storing Receipts and Logs in the DB - we won't be doing that in Erigon 3
// do not regen, more context: https://github.com/erigontech/erigon/pull/10105#pullrequestreview-2027423601
// go:generate codecgen -o receipt_codecgen_gen.go -r "^Receipts$|^Receipt$|^Logs$|^Log$" -st "codec" -j=false -nx=true -ta=true -oe=false -d 2 receipt.go log.go

//go:generate gencodec -type Receipt -field-override receiptMarshaling -out gen_receipt_json.go
//go:generate codecgen -o receipt_codecgen_gen.go -r "^Receipts$|^Receipt$|^Logs$|^Log$" -st "codec" -j=false -nx=true -ta=true -oe=false -d 2 receipt.go log.go

var (
	receiptStatusFailedRLP     = []byte{}
	receiptStatusSuccessfulRLP = []byte{0x01}
)

const (
	// ReceiptStatusFailed is the status code of a transaction if execution failed.
	ReceiptStatusFailed = uint64(0)

	// ReceiptStatusSuccessful is the status code of a transaction if execution succeeded.
	ReceiptStatusSuccessful = uint64(1)

	// The version number for post-canyon deposit receipts.
	CanyonDepositReceiptVersion = uint64(1)
)

// Receipt represents the results of a transaction.
// DESCRIBED: docs/programmers_guide/guide.md#organising-ethereum-state-into-a-merkle-tree
type Receipt struct {
	// Consensus fields: These fields are defined by the Yellow Paper
	Type              uint8  `json:"type,omitempty"`
	PostState         []byte `json:"root" codec:"1"`
	Status            uint64 `json:"status" codec:"2"`
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed" gencodec:"required" codec:"3"`
	Bloom             Bloom  `json:"logsBloom"         gencodec:"required" codec:"-"`
	Logs              Logs   `json:"logs"              gencodec:"required" codec:"-"`

	// Implementation fields: These fields are added by geth when processing a transaction.
	// They are stored in the chain database.
	TxHash          libcommon.Hash    `json:"transactionHash" gencodec:"required" codec:"-"`
	ContractAddress libcommon.Address `json:"contractAddress" codec:"-"`
	GasUsed         uint64            `json:"gasUsed" gencodec:"required" codec:"-"`

	// DepositNonce was introduced in Regolith to store the actual nonce used by deposit transactions
	// The state transition process ensures this is only set for Regolith deposit transactions.
	DepositNonce *uint64 `json:"depositNonce,omitempty"`
	// The position of DepositNonce variable must NOT be changed. If changed, cbor decoding will fail
	// for the data following previous struct and leading to decoding error(triggering backward imcompatibility).

	// Further fields when added must be appended after the last variable. Watch out for cbor.

	// Inclusion information: These fields provide information about the inclusion of the
	// transaction corresponding to this receipt.
	BlockHash        libcommon.Hash `json:"blockHash,omitempty" codec:"-"`
	BlockNumber      *big.Int       `json:"blockNumber,omitempty" codec:"-"`
	TransactionIndex uint           `json:"transactionIndex" codec:"-"`

	// OVM legacy: extend receipts with their L1 price (if a rollup tx)
	L1GasPrice *big.Int   `json:"l1GasPrice,omitempty"`
	L1GasUsed  *big.Int   `json:"l1GasUsed,omitempty"`   // Present from pre-bedrock, deprecated as of Fjord
	L1Fee      *big.Int   `json:"l1Fee,omitempty"`       // Present from pre-bedrock
	FeeScalar  *big.Float `json:"l1FeeScalar,omitempty"` // Present from pre-bedrock to Ecotone. Nil after Ecotone

	// DepositReceiptVersion was introduced in Canyon to indicate an update to how receipt hashes
	// should be computed when set. The state transition process ensures this is only set for
	// post-Canyon deposit transactions.
	DepositReceiptVersion *uint64 `json:"depositReceiptVersion,omitempty"`

	// Fee scalars were introduced after the Ecotone hardfork
	L1BaseFeeScalar     *uint64  `json:"l1BaseFeeScalar,omitempty"`     // Always nil prior to the Ecotone hardfork
	L1BlobBaseFeeScalar *uint64  `json:"l1BlobBaseFeeScalar,omitempty"` // Always nil prior to the Ecotone hardfork
	L1BlobBaseFee       *big.Int `json:"l1BlobBaseFee,omitempty"`       // Always nil prior to the Ecotone hardfork
}

type receiptMarshaling struct {
	Type              hexutil.Uint64
	PostState         hexutility.Bytes
	Status            hexutil.Uint64
	CumulativeGasUsed hexutil.Uint64
	GasUsed           hexutil.Uint64
	BlockNumber       *hexutil.Big
	TransactionIndex  hexutil.Uint

	// Optimism
	L1GasPrice            *hexutil.Big
	L1GasUsed             *hexutil.Big
	L1Fee                 *hexutil.Big
	FeeScalar             *big.Float
	DepositNonce          *hexutil.Uint64
	DepositReceiptVersion *hexutil.Uint64
	L1BaseFeeScalar       *hexutil.Uint64
	L1BlobBaseFee         *hexutil.Big
	L1BlobBaseFeeScalar   *hexutil.Uint64
}

// receiptRLP is the consensus encoding of a receipt.
type receiptRLP struct {
	PostStateOrStatus []byte
	CumulativeGasUsed uint64
	Bloom             Bloom
	Logs              []*Log
}

type depositReceiptRlp struct {
	PostStateOrStatus []byte
	CumulativeGasUsed uint64
	Bloom             Bloom
	Logs              []*Log
	// DepositNonce was introduced in Regolith to store the actual nonce used by deposit transactions.
	// Must be nil for any transactions prior to Regolith or that aren't deposit transactions.
	DepositNonce *uint64 `rlp:"optional"`
	// Receipt hash post-Regolith but pre-Canyon inadvertently did not include the above
	// DepositNonce. Post Canyon, receipts will have a non-empty DepositReceiptVersion indicating
	// which post-Canyon receipt hash function to invoke.
	DepositReceiptVersion *uint64 `rlp:"optional"`
}

// storedReceiptRLP is the storage encoding of a receipt.
type storedReceiptRLP struct {
	PostStateOrStatus []byte
	CumulativeGasUsed uint64
	Logs              []*LogForStorage
	// DepositNonce was introduced in Regolith to store the actual nonce used by deposit transactions.
	// Must be nil for any transactions prior to Regolith or that aren't deposit transactions.
	DepositNonce *uint64 `rlp:"optional"`
	// Receipt hash post-Regolith but pre-Canyon inadvertently did not include the above
	// DepositNonce. Post Canyon, receipts will have a non-empty DepositReceiptVersion indicating
	// which post-Canyon receipt hash function to invoke.
	DepositReceiptVersion *uint64 `rlp:"optional"`
}

// v4StoredReceiptRLP is the storage encoding of a receipt used in database version 4.
type v4StoredReceiptRLP struct {
	PostStateOrStatus []byte
	CumulativeGasUsed uint64
	TxHash            libcommon.Hash
	ContractAddress   libcommon.Address
	Logs              []*LogForStorage
	GasUsed           uint64
}

// v3StoredReceiptRLP is the original storage encoding of a receipt including some unnecessary fields.
type v3StoredReceiptRLP struct {
	PostStateOrStatus []byte
	CumulativeGasUsed uint64
	//Bloom             Bloom
	//TxHash            libcommon.Hash
	ContractAddress libcommon.Address
	Logs            []*LogForStorage
	GasUsed         uint64
}

// NewReceipt creates a barebone transaction receipt, copying the init fields.
// Deprecated: create receipts using a struct literal instead.
func NewReceipt(failed bool, cumulativeGasUsed uint64) *Receipt {
	r := &Receipt{
		Type:              LegacyTxType,
		CumulativeGasUsed: cumulativeGasUsed,
	}
	if failed {
		r.Status = ReceiptStatusFailed
	} else {
		r.Status = ReceiptStatusSuccessful
	}
	return r
}

// EncodeRLP implements rlp.Encoder, and flattens the consensus fields of a receipt
// into an RLP stream. If no post state is present, byzantium fork is assumed.
func (r Receipt) EncodeRLP(w io.Writer) error {
	data := &receiptRLP{r.statusEncoding(), r.CumulativeGasUsed, r.Bloom, r.Logs}
	if r.Type == LegacyTxType {
		return rlp.Encode(w, data)
	}
	buf := new(bytes.Buffer)
	buf.WriteByte(r.Type)
	if r.Type == DepositTxType {
		withNonceAndReceiptVersion := &depositReceiptRlp{r.statusEncoding(), r.CumulativeGasUsed, r.Bloom, r.Logs, r.DepositNonce, r.DepositReceiptVersion}
		if err := rlp.Encode(buf, withNonceAndReceiptVersion); err != nil {
			return err
		}
	} else {
		if err := rlp.Encode(buf, data); err != nil {
			return err
		}
	}
	return rlp.Encode(w, buf.Bytes())
}

// encodeTyped writes the canonical encoding of a typed receipt to w.
func (r *Receipt) encodeTyped(data *receiptRLP, w *bytes.Buffer) error {
	w.WriteByte(r.Type)
	switch r.Type {
	case DepositTxType:
		withNonceAndReceiptVersion := depositReceiptRlp{data.PostStateOrStatus, data.CumulativeGasUsed, data.Bloom, data.Logs, r.DepositNonce, r.DepositReceiptVersion}
		return rlp.Encode(w, withNonceAndReceiptVersion)
	default:
		return rlp.Encode(w, data)
	}
}
func (r *Receipt) decodePayload(s *rlp.Stream) error {
	_, err := s.List()
	if err != nil {
		return err
	}
	var b []byte
	if b, err = s.Bytes(); err != nil {
		return fmt.Errorf("read PostStateOrStatus: %w", err)
	}
	r.setStatus(b)
	if r.CumulativeGasUsed, err = s.Uint(); err != nil {
		return fmt.Errorf("read CumulativeGasUsed: %w", err)
	}
	if b, err = s.Bytes(); err != nil {
		return fmt.Errorf("read Bloom: %w", err)
	}
	if len(b) != 256 {
		return fmt.Errorf("wrong size for Bloom: %d", len(b))
	}
	copy(r.Bloom[:], b)
	// decode logs
	if _, err = s.List(); err != nil {
		return fmt.Errorf("open Logs: %w", err)
	}
	if r.Logs != nil && len(r.Logs) > 0 {
		r.Logs = r.Logs[:0]
	}
	for _, err = s.List(); err == nil; _, err = s.List() {
		r.Logs = append(r.Logs, &Log{})
		log := r.Logs[len(r.Logs)-1]
		if b, err = s.Bytes(); err != nil {
			return fmt.Errorf("read Address: %w", err)
		}
		if len(b) != 20 {
			return fmt.Errorf("wrong size for Log address: %d", len(b))
		}
		copy(log.Address[:], b)
		if _, err = s.List(); err != nil {
			return fmt.Errorf("open Topics: %w", err)
		}
		for b, err = s.Bytes(); err == nil; b, err = s.Bytes() {
			log.Topics = append(log.Topics, libcommon.Hash{})
			if len(b) != 32 {
				return fmt.Errorf("wrong size for Topic: %d", len(b))
			}
			copy(log.Topics[len(log.Topics)-1][:], b)
		}
		if !errors.Is(err, rlp.EOL) {
			return fmt.Errorf("read Topic: %w", err)
		}
		// end of Topics list
		if err = s.ListEnd(); err != nil {
			return fmt.Errorf("close Topics: %w", err)
		}
		if log.Data, err = s.Bytes(); err != nil {
			return fmt.Errorf("read Data: %w", err)
		}
		// end of Log
		if err = s.ListEnd(); err != nil {
			return fmt.Errorf("close Log: %w", err)
		}
	}
	if !errors.Is(err, rlp.EOL) {
		return fmt.Errorf("open Log: %w", err)
	}
	if err = s.ListEnd(); err != nil {
		return fmt.Errorf("close Logs: %w", err)
	}
	if r.Type == DepositTxType {
		depositNonce, err := s.Uint()
		if err != nil {
			if !errors.Is(err, rlp.EOL) {
				return fmt.Errorf("read DepositNonce: %w", err)
			}
			return nil
		} else {
			r.DepositNonce = &depositNonce
		}
		depositReceiptVersion, err := s.Uint()
		if err != nil {
			if !errors.Is(err, rlp.EOL) {
				return fmt.Errorf("read DepositReceiptVersion: %w", err)
			}
			return nil
		} else {
			r.DepositReceiptVersion = &depositReceiptVersion
		}
	}
	if err := s.ListEnd(); err != nil {
		return fmt.Errorf("close receipt payload: %w", err)
	}
	return nil
}

// DecodeRLP implements rlp.Decoder, and loads the consensus fields of a receipt
// from an RLP stream.
func (r *Receipt) DecodeRLP(s *rlp.Stream) error {
	kind, size, err := s.Kind()
	if err != nil {
		return err
	}
	switch kind {
	case rlp.List:
		// It's a legacy receipt.
		if err := r.decodePayload(s); err != nil {
			return err
		}
		r.Type = LegacyTxType
	case rlp.String:
		// It's an EIP-2718 typed tx receipt.
		s.NewList(size) // Hack - convert String (envelope) into List
		var b []byte
		if b, err = s.Bytes(); err != nil {
			return fmt.Errorf("read TxType: %w", err)
		}
		if len(b) != 1 {
			return fmt.Errorf("%w, got %d bytes", rlp.ErrWrongTxTypePrefix, len(b))
		}
		r.Type = b[0]
		switch r.Type {
		case AccessListTxType, DynamicFeeTxType, DepositTxType, BlobTxType, SetCodeTxType:
			if err := r.decodePayload(s); err != nil {
				return err
			}
		default:
			return ErrTxTypeNotSupported
		}
		if err = s.ListEnd(); err != nil {
			return err
		}
	default:
		return rlp.ErrExpectedList
	}
	return nil
}

func (r *Receipt) setStatus(postStateOrStatus []byte) error {
	switch {
	case bytes.Equal(postStateOrStatus, receiptStatusSuccessfulRLP):
		r.Status = ReceiptStatusSuccessful
	case bytes.Equal(postStateOrStatus, receiptStatusFailedRLP):
		r.Status = ReceiptStatusFailed
	case len(postStateOrStatus) == len(libcommon.Hash{}):
		r.PostState = postStateOrStatus
	default:
		return fmt.Errorf("invalid receipt status %x", postStateOrStatus)
	}
	return nil
}

func (r *Receipt) statusEncoding() []byte {
	if len(r.PostState) == 0 {
		if r.Status == ReceiptStatusFailed {
			return receiptStatusFailedRLP
		}
		return receiptStatusSuccessfulRLP
	}
	return r.PostState
}

// Copy creates a deep copy of the Receipt.
func (r *Receipt) Copy() *Receipt {
	postState := make([]byte, len(r.PostState))
	copy(postState, r.PostState)

	bloom := BytesToBloom(r.Bloom.Bytes())

	logs := make(Logs, 0, len(r.Logs))
	for _, log := range r.Logs {
		logs = append(logs, log.Copy())
	}

	txHash := libcommon.BytesToHash(r.TxHash.Bytes())
	contractAddress := libcommon.BytesToAddress(r.ContractAddress.Bytes())
	blockHash := libcommon.BytesToHash(r.BlockHash.Bytes())
	var blockNumber *big.Int
	if r.BlockNumber != nil {
		blockNumber = big.NewInt(0).Set(r.BlockNumber)
	}

	return &Receipt{
		Type:                  r.Type,
		PostState:             postState,
		Status:                r.Status,
		CumulativeGasUsed:     r.CumulativeGasUsed,
		Bloom:                 bloom,
		Logs:                  logs,
		TxHash:                txHash,
		ContractAddress:       contractAddress,
		GasUsed:               r.GasUsed,
		BlockHash:             blockHash,
		BlockNumber:           blockNumber,
		TransactionIndex:      r.TransactionIndex,
		DepositNonce:          r.DepositNonce,
		DepositReceiptVersion: r.DepositReceiptVersion,
	}
}

type ReceiptsForStorage []*ReceiptForStorage

// ReceiptForStorage is a wrapper around a Receipt that flattens and parses the
// entire content of a receipt, as opposed to only the consensus fields originally.
type ReceiptForStorage Receipt

// EncodeRLP implements rlp.Encoder, and flattens all content fields of a receipt
// into an RLP stream.
func (r *ReceiptForStorage) EncodeRLP(w io.Writer) error {
	enc := &storedReceiptRLP{
		PostStateOrStatus:     (*Receipt)(r).statusEncoding(),
		CumulativeGasUsed:     r.CumulativeGasUsed,
		Logs:                  make([]*LogForStorage, len(r.Logs)),
		DepositNonce:          r.DepositNonce,
		DepositReceiptVersion: r.DepositReceiptVersion,
	}
	for i, log := range r.Logs {
		enc.Logs[i] = (*LogForStorage)(log)
	}
	return rlp.Encode(w, enc)
}

// DecodeRLP implements rlp.Decoder, and loads both consensus and implementation
// fields of a receipt from an RLP stream.
func (r *ReceiptForStorage) DecodeRLP(s *rlp.Stream) error {
	// Retrieve the entire receipt blob as we need to try multiple decoders
	blob, err := s.Raw()
	if err != nil {
		return err
	}
	// Try decoding from the newest format for future proofness, then the older one
	// for old nodes that just upgraded. V4 was an intermediate unreleased format so
	// we do need to decode it, but it's not common (try last).
	if err := decodeStoredReceiptRLP(r, blob); err == nil {
		return nil
	}
	if err := decodeV3StoredReceiptRLP(r, blob); err == nil {
		return nil
	}
	return decodeV4StoredReceiptRLP(r, blob)
}

func decodeStoredReceiptRLP(r *ReceiptForStorage, blob []byte) error {
	var stored storedReceiptRLP
	if err := rlp.DecodeBytes(blob, &stored); err != nil {
		return err
	}
	if err := (*Receipt)(r).setStatus(stored.PostStateOrStatus); err != nil {
		return err
	}
	r.CumulativeGasUsed = stored.CumulativeGasUsed
	r.Logs = make([]*Log, len(stored.Logs))
	for i, log := range stored.Logs {
		r.Logs[i] = (*Log)(log)
	}
	//r.Bloom = CreateBloom(Receipts{(*Receipt)(r)})
	if stored.DepositNonce != nil {
		r.DepositNonce = stored.DepositNonce
	}

	if stored.DepositReceiptVersion != nil {
		r.DepositReceiptVersion = stored.DepositReceiptVersion
	}
	return nil
}

func decodeV4StoredReceiptRLP(r *ReceiptForStorage, blob []byte) error {
	var stored v4StoredReceiptRLP
	if err := rlp.DecodeBytes(blob, &stored); err != nil {
		return err
	}
	if err := (*Receipt)(r).setStatus(stored.PostStateOrStatus); err != nil {
		return err
	}
	r.CumulativeGasUsed = stored.CumulativeGasUsed
	r.TxHash = stored.TxHash
	r.ContractAddress = stored.ContractAddress
	r.GasUsed = stored.GasUsed
	r.Logs = make([]*Log, len(stored.Logs))
	for i, log := range stored.Logs {
		r.Logs[i] = (*Log)(log)
	}
	//r.Bloom = CreateBloom(Receipts{(*Receipt)(r)})

	return nil
}

func decodeV3StoredReceiptRLP(r *ReceiptForStorage, blob []byte) error {
	var stored v3StoredReceiptRLP
	if err := rlp.DecodeBytes(blob, &stored); err != nil {
		return err
	}
	if err := (*Receipt)(r).setStatus(stored.PostStateOrStatus); err != nil {
		return err
	}
	r.CumulativeGasUsed = stored.CumulativeGasUsed
	r.ContractAddress = stored.ContractAddress
	r.GasUsed = stored.GasUsed
	r.Logs = make([]*Log, len(stored.Logs))
	for i, log := range stored.Logs {
		r.Logs[i] = (*Log)(log)
	}
	return nil
}

// Receipts implements DerivableList for receipts.
type Receipts []*Receipt

// Len returns the number of receipts in this list.
func (rs Receipts) Len() int { return len(rs) }

// EncodeIndex encodes the i'th receipt to w.
// During post-regolith and pre-Canyon, DepositNonce was not included when encoding for hashing.
// Canyon adds DepositReceiptVersion to preserve backwards compatibility for pre-Canyon, and
// for correct receipt-root hash computation.
func (rs Receipts) EncodeIndex(i int, w *bytes.Buffer) {
	r := rs[i]
	data := &receiptRLP{r.statusEncoding(), r.CumulativeGasUsed, r.Bloom, r.Logs}
	switch r.Type {
	case LegacyTxType:
		if err := rlp.Encode(w, data); err != nil {
			panic(err)
		}
	case AccessListTxType:
		//nolint:errcheck
		w.WriteByte(AccessListTxType)
		if err := rlp.Encode(w, data); err != nil {
			panic(err)
		}
	case DynamicFeeTxType:
		w.WriteByte(DynamicFeeTxType)
		if err := rlp.Encode(w, data); err != nil {
			panic(err)
		}
	case DepositTxType:
		w.WriteByte(DepositTxType)
		if r.DepositReceiptVersion != nil {
			// post-canyon receipt hash computation update
			depositData := &depositReceiptRlp{data.PostStateOrStatus, data.CumulativeGasUsed, r.Bloom, r.Logs, r.DepositNonce, r.DepositReceiptVersion}
			if err := rlp.Encode(w, depositData); err != nil {
				panic(err)
			}
		} else {
			if err := rlp.Encode(w, data); err != nil {
				panic(err)
			}
		}
	case BlobTxType:
		w.WriteByte(BlobTxType)
		if err := rlp.Encode(w, data); err != nil {
			panic(err)
		}
	case SetCodeTxType:
		w.WriteByte(SetCodeTxType)
		if err := rlp.Encode(w, data); err != nil {
			panic(err)
		}
	default:
		// For unsupported types, write nothing. Since this is for
		// DeriveSha, the error will be caught matching the derived hash
		// to the block.
	}
}

// DeriveFields fills the receipts with their computed fields based on consensus
// data and contextual infos like containing block and transactions.
func (r Receipts) DeriveFields(config *chain.Config, hash libcommon.Hash, number uint64, time uint64, txs Transactions, senders []libcommon.Address) error {
	logIndex := uint(0) // logIdx is unique within the block and starts from 0
	if len(txs) != len(r) {
		return fmt.Errorf("transaction and receipt count mismatch, tx count = %d, receipts count = %d", len(txs), len(r))
	}
	if len(senders) != len(txs) {
		return fmt.Errorf("transaction and senders count mismatch, tx count = %d, senders count = %d", len(txs), len(senders))
	}

	blockNumber := new(big.Int).SetUint64(number)
	for i := 0; i < len(r); i++ {
		// The transaction type and hash can be retrieved from the transaction itself
		r[i].Type = txs[i].Type()
		r[i].TxHash = txs[i].Hash()

		// block location fields
		r[i].BlockHash = hash
		r[i].BlockNumber = blockNumber
		r[i].TransactionIndex = uint(i)

		// The contract address can be derived from the transaction itself
		if txs[i].GetTo() == nil {
			// If one wants to deploy a contract, one needs to send a transaction that does not have `To` field
			// and then the address of the contract one is creating this way will depend on the `tx.From`
			// and the nonce of the creating account (which is `tx.From`).
			nonce := txs[i].GetNonce()
			if r[i].DepositNonce != nil {
				nonce = *r[i].DepositNonce
			}
			r[i].ContractAddress = crypto.CreateAddress(senders[i], nonce)
		}
		// The used gas can be calculated based on previous r
		if i == 0 {
			r[i].GasUsed = r[i].CumulativeGasUsed
		} else {
			r[i].GasUsed = r[i].CumulativeGasUsed - r[i-1].CumulativeGasUsed
		}
		// The derived log fields can simply be set from the block and transaction
		for j := 0; j < len(r[i].Logs); j++ {
			r[i].Logs[j].BlockNumber = number
			r[i].Logs[j].BlockHash = hash
			r[i].Logs[j].TxHash = r[i].TxHash
			r[i].Logs[j].TxIndex = uint(i)
			r[i].Logs[j].Index = logIndex
			logIndex++
		}
	}
	if config.IsOptimismBedrock(number) && len(txs) >= 2 { // need at least an info tx and a non-info tx
		gasParams, err := opstack.ExtractL1GasParams(config, time, txs[0].GetData())
		if err != nil {
			return err
		}
		for i := 0; i < len(r); i++ {
			if txs[i].Type() == DepositTxType {
				continue
			}
			r[i].L1GasPrice = gasParams.L1BaseFee.ToBig()
			l1Fee, l1GasUsed := gasParams.CostFunc(txs[i].RollupCostData())
			r[i].L1Fee = l1Fee.ToBig()
			r[i].L1GasUsed = l1GasUsed.ToBig()
			if gasParams.FeeScalar != nil {
				r[i].FeeScalar = gasParams.FeeScalar
			}
			if gasParams.L1BaseFeeScalar != nil {
				l1BaseFeeScalar := gasParams.L1BaseFeeScalar.Uint64()
				r[i].L1BaseFeeScalar = &l1BaseFeeScalar
			}
			if gasParams.L1BlobBaseFee != nil {
				l1BlobBaseFee := gasParams.L1BlobBaseFee
				r[i].L1BlobBaseFee = l1BlobBaseFee.ToBig()
			}
			if gasParams.L1BlobBaseFeeScalar != nil {
				l1BlobBaseFeeScalar := gasParams.L1BlobBaseFeeScalar.Uint64()
				r[i].L1BlobBaseFeeScalar = &l1BlobBaseFeeScalar
			}
		}
	}
	return nil
}
