package rawdbhelpers

import (
	"encoding/binary"

	"github.com/erigontech/erigon-lib/config3"
	"github.com/erigontech/erigon-lib/kv"
)

func IdxStepsCountV3(tx kv.Tx) float64 {
	fst, _ := kv.FirstKey(tx, kv.TblTracesToKeys)
	lst, _ := kv.LastKey(tx, kv.TblTracesToKeys)
	if len(fst) > 0 && len(lst) > 0 {
		fstTxNum := binary.BigEndian.Uint64(fst)
		lstTxNum := binary.BigEndian.Uint64(lst)

		return float64(lstTxNum-fstTxNum) / float64(config3.HistoryV3AggregationStep)
	}
	return 0
}
