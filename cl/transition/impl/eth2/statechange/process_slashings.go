package statechange

import (
	"github.com/erigontech/erigon/cl/abstract"
	"github.com/erigontech/erigon/cl/clparams"
	"github.com/erigontech/erigon/cl/cltypes/solid"
	"github.com/erigontech/erigon/cl/phase1/core/state"
)

func processSlashings(s abstract.BeaconState, slashingMultiplier uint64) error {
	// Get the current epoch
	epoch := state.Epoch(s)
	// Get the total active balance
	totalBalance := s.GetTotalActiveBalance()
	// Calculate the total slashing amount
	// by summing all slashings and multiplying by the provided multiplier
	slashing := state.GetTotalSlashingAmount(s) * slashingMultiplier
	// Adjust the total slashing amount to be no greater than the total active balance
	if totalBalance < slashing {
		slashing = totalBalance
	}
	beaconConfig := s.BeaconConfig()
	// Apply penalties to validators who have been slashed and reached the withdrawable epoch
	var err error
	s.ForEachValidator(func(validator solid.Validator, i, total int) bool {
		if !validator.Slashed() || epoch+beaconConfig.EpochsPerSlashingsVector/2 != validator.WithdrawableEpoch() {
			return true
		}
		// Get the effective balance increment
		increment := beaconConfig.EffectiveBalanceIncrement
		// Calculate the penalty numerator by multiplying the validator's effective balance by the total slashing amount
		penaltyNumerator := validator.EffectiveBalance() / increment * slashing
		// Calculate the penalty by dividing the penalty numerator by the total balance and multiplying by the increment
		penalty := penaltyNumerator / totalBalance * increment
		// Decrease the validator's balance by the calculated penalty
		if err = state.DecreaseBalance(s, uint64(i), penalty); err != nil {
			return false
		}
		return true
	})
	if err != nil {
		return err
	}
	return nil
}

func ProcessSlashings(state abstract.BeaconState) error {
	// Depending on the version of the state, use different multipliers
	switch state.Version() {
	case clparams.Phase0Version:
		return processSlashings(state, state.BeaconConfig().ProportionalSlashingMultiplier)
	case clparams.AltairVersion:
		return processSlashings(state, state.BeaconConfig().ProportionalSlashingMultiplierAltair)
	default:
		return processSlashings(state, state.BeaconConfig().ProportionalSlashingMultiplierBellatrix)
	}
}
