package misc

import (
	"github.com/erigontech/erigon-lib/log/v3"

	libcommon "github.com/erigontech/erigon-lib/common"
	"github.com/erigontech/erigon/consensus"
	"github.com/erigontech/erigon/params"
)

func ApplyBeaconRootEip4788(parentBeaconBlockRoot *libcommon.Hash, syscall consensus.SystemCall) {
	if parentBeaconBlockRoot == nil {
		log.Warn("Skipping EIP-4788 as there is no parentBeaconBlockRoot")
		return
	}
	_, err := syscall(params.BeaconRootsAddress, parentBeaconBlockRoot.Bytes())
	if err != nil {
		log.Warn("Failed to call beacon roots contract", "err", err)
	}
}
