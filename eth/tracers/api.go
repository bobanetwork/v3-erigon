package tracers

import (
	"encoding/json"

	"github.com/erigontech/erigon-lib/common/hexutil"
	"github.com/erigontech/erigon/eth/tracers/logger"
	"github.com/erigontech/erigon/turbo/adapter/ethapi"
)

// TraceConfig holds extra parameters to trace functions.
type TraceConfig struct {
	*logger.LogConfig
	Tracer         *string
	TracerConfig   *json.RawMessage
	Timeout        *string
	Reexec         *uint64
	NoRefunds      *bool // Turns off gas refunds when tracing
	StateOverrides *ethapi.StateOverrides

	BorTraceEnabled *bool
	TxIndex         *hexutil.Uint
}
