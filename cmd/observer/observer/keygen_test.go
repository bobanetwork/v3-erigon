package observer

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/erigontech/erigon-lib/log/v3"
	"github.com/erigontech/erigon/crypto"
	"github.com/stretchr/testify/assert"
)

func TestKeygen(t *testing.T) {
	targetKeyPair, err := crypto.GenerateKey()
	assert.NotNil(t, targetKeyPair)
	assert.Nil(t, err)

	targetKey := &targetKeyPair.PublicKey
	keys := keygen(context.Background(), targetKey, 50*time.Millisecond, uint(runtime.GOMAXPROCS(-1)), log.Root())

	assert.NotNil(t, keys)
	assert.GreaterOrEqual(t, len(keys), 4)
}
