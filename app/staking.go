package app

import (
	"time"

	sdk "cosmossdk.io/math"
)

var (
	// StakingParameters defines the parameters for staking module
	StakingParameters = StakingParams{
		UnbondingTime:     time.Duration(21) * 24 * time.Hour, // 21 days
		MaxValidators:     3,
		MaxEntries:        7,
		HistoricalEntries: 10000,
		MinCommissionRate: sdk.LegacyNewDecWithPrec(5, 2), // 5%
	}
)

type StakingParams struct {
	UnbondingTime     time.Duration
	MaxValidators     uint32
	MaxEntries        uint32
	HistoricalEntries uint32
	MinCommissionRate sdk.LegacyDec
}
