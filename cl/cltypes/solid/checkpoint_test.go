package solid_test

import (
	"testing"

	libcommon "github.com/erigontech/erigon-lib/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/erigontech/erigon/cl/cltypes/solid"
)

var testCheckpoint = solid.NewCheckpointFromParameters(libcommon.HexToHash("0x3"), 69)

var expectedTestCheckpointMarshalled = libcommon.Hex2Bytes("45000000000000000000000000000000000000000000000000000000000000000000000000000003")
var expectedTestCheckpointRoot = libcommon.Hex2Bytes("be8567f9fdae831b10720823dbcf0e3680e61d6a2a27d85ca00f6c15a7bbb1ea")

func TestCheckpointMarshalUnmarmashal(t *testing.T) {
	marshalled, err := testCheckpoint.EncodeSSZ(nil)
	require.NoError(t, err)
	assert.Equal(t, marshalled, expectedTestCheckpointMarshalled)
	checkpoint := solid.NewCheckpoint()
	require.NoError(t, checkpoint.DecodeSSZ(marshalled, 0))
	require.Equal(t, checkpoint, testCheckpoint)
}

func TestCheckpointHashTreeRoot(t *testing.T) {
	root, err := testCheckpoint.HashSSZ()
	require.NoError(t, err)
	assert.Equal(t, root[:], expectedTestCheckpointRoot)
}
