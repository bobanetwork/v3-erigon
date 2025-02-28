package state

import (
	"github.com/erigontech/erigon-lib/metrics"
	"github.com/erigontech/erigon-lib/types/clonable"
)

func (b *CachingBeaconState) EncodeSSZ(buf []byte) ([]byte, error) {
	h := metrics.NewHistTimer("encode_ssz_beacon_state_dur")
	bts, err := b.BeaconState.EncodeSSZ(buf)
	if err != nil {
		return nil, err
	}
	h.PutSince()
	sz := metrics.NewHistTimer("encode_ssz_beacon_state_size")
	sz.Observe(float64(len(bts)))
	return bts, err
}

func (b *CachingBeaconState) DecodeSSZ(buf []byte, version int) error {
	h := metrics.NewHistTimer("decode_ssz_beacon_state_dur")
	if err := b.BeaconState.DecodeSSZ(buf, version); err != nil {
		return err
	}
	sz := metrics.NewHistTimer("decode_ssz_beacon_state_size")
	sz.Observe(float64(len(buf)))
	h.PutSince()
	return b.InitBeaconState()
}

// SSZ size of the Beacon State
func (b *CachingBeaconState) EncodingSizeSSZ() (size int) {
	sz := b.BeaconState.EncodingSizeSSZ()
	return sz
}

func (b *CachingBeaconState) Clone() clonable.Clonable {
	return New(b.BeaconConfig())
}
