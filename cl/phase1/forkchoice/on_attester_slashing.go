package forkchoice

import (
	"fmt"

	"github.com/Giulio2002/bls"
	"github.com/erigontech/erigon/cl/cltypes/solid"
	"github.com/erigontech/erigon/cl/fork"
	"github.com/erigontech/erigon/cl/phase1/core/state"
	"github.com/erigontech/erigon/cl/pool"

	"github.com/erigontech/erigon/cl/cltypes"
)

func (f *ForkChoiceStore) OnAttesterSlashing(attesterSlashing *cltypes.AttesterSlashing, test bool) error {
	if f.operationsPool.AttesterSlashingsPool.Has(pool.ComputeKeyForAttesterSlashing(attesterSlashing)) {
		return nil
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	// Check if this attestation is even slashable.
	attestation1 := attesterSlashing.Attestation_1
	attestation2 := attesterSlashing.Attestation_2
	if !cltypes.IsSlashableAttestationData(attestation1.Data, attestation2.Data) {
		return fmt.Errorf("attestation data is not slashable")
	}
	var err error
	s := f.syncedDataManager.HeadState()
	if s == nil {
		// Retrieve justified state
		s, err = f.forkGraph.GetState(f.justifiedCheckpoint.Load().(solid.Checkpoint).BlockRoot(), false)
		if err != nil {
			return err
		}
	}
	if s == nil {
		return fmt.Errorf("no state accessible")
	}
	attestation1PublicKeys, err := getIndexedAttestationPublicKeys(s, attestation1)
	if err != nil {
		return err
	}
	attestation2PublicKeys, err := getIndexedAttestationPublicKeys(s, attestation2)
	if err != nil {
		return err
	}
	domain1, err := s.GetDomain(s.BeaconConfig().DomainBeaconAttester, attestation1.Data.Target().Epoch())
	if err != nil {
		return fmt.Errorf("unable to get the domain: %v", err)
	}
	domain2, err := s.GetDomain(s.BeaconConfig().DomainBeaconAttester, attestation2.Data.Target().Epoch())
	if err != nil {
		return fmt.Errorf("unable to get the domain: %v", err)
	}

	if !test {
		// Verify validity of slashings (1)
		signingRoot, err := fork.ComputeSigningRoot(attestation1.Data, domain1)
		if err != nil {
			return fmt.Errorf("unable to get signing root: %v", err)
		}

		valid, err := bls.VerifyAggregate(attestation1.Signature[:], signingRoot[:], attestation1PublicKeys)
		if err != nil {
			return fmt.Errorf("error while validating signature: %v", err)
		}
		if !valid {
			return fmt.Errorf("invalid aggregate signature")
		}
		// Verify validity of slashings (2)
		signingRoot, err = fork.ComputeSigningRoot(attestation2.Data, domain2)
		if err != nil {
			return fmt.Errorf("unable to get signing root: %v", err)
		}

		valid, err = bls.VerifyAggregate(attestation2.Signature[:], signingRoot[:], attestation2PublicKeys)
		if err != nil {
			return fmt.Errorf("error while validating signature: %v", err)
		}
		if !valid {
			return fmt.Errorf("invalid aggregate signature")
		}
	}

	var anySlashed bool
	for _, index := range solid.IntersectionOfSortedSets(attestation1.AttestingIndices, attestation2.AttestingIndices) {
		f.setUnequivocating(index)
		if !anySlashed {
			v, err := s.ValidatorForValidatorIndex(int(index))
			if err != nil {
				return fmt.Errorf("unable to retrieve state: %v", err)
			}
			if v.IsSlashable(state.Epoch(s)) {
				anySlashed = true
			}
		}
	}
	if anySlashed {
		f.operationsPool.AttesterSlashingsPool.Insert(pool.ComputeKeyForAttesterSlashing(attesterSlashing), attesterSlashing)
	}
	return nil
}

func getIndexedAttestationPublicKeys(b *state.CachingBeaconState, att *cltypes.IndexedAttestation) ([][]byte, error) {
	inds := att.AttestingIndices
	if inds.Length() == 0 || !solid.IsUint64SortedSet(inds) {
		return nil, fmt.Errorf("isValidIndexedAttestation: attesting indices are not sorted or are null")
	}
	pks := make([][]byte, 0, inds.Length())
	if err := solid.RangeErr[uint64](inds, func(_ int, v uint64, _ int) error {
		val, err := b.ValidatorForValidatorIndex(int(v))
		if err != nil {
			return err
		}
		pk := val.PublicKey()
		pks = append(pks, pk[:])
		return nil
	}); err != nil {
		return nil, err
	}
	return pks, nil
}
