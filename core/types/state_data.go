package types

import (
	libcommon "github.com/erigontech/erigon-lib/common"
)

// StateSyncData represents state received from Ethereum Blockchain
type StateSyncData struct {
	ID       uint64
	Contract libcommon.Address
	Data     string
	TxHash   libcommon.Hash
}
