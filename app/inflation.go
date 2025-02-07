package app

import (
	sdk "cosmossdk.io/math"
)

var (
	// InflationParameters defines the parameters for token inflation
	InflationParameters = InflationParams{
		InitialInflation:  sdk.LegacyNewDecWithPrec(143, 3), // 14.3% first year
		ReductionFactor:   sdk.LegacyNewDecWithPrec(97, 2),  // Reduce by 3% each year
		BlocksPerYear:     uint64(6_311_520),                // ~5 second blocks
		TargetBondedRatio: sdk.LegacyNewDecWithPrec(67, 2),  // 67% target bonded ratio
		InflationDistribution: InflationDistributionParams{
			StakingRewards: sdk.LegacyNewDecWithPrec(533333334, 9), // ~53.33%
			CommunityPool:  sdk.LegacyNewDecWithPrec(466666666, 9), // ~46.67%
		},
	}
)

type InflationParams struct {
	InitialInflation      sdk.LegacyDec
	ReductionFactor       sdk.LegacyDec
	BlocksPerYear         uint64
	TargetBondedRatio     sdk.LegacyDec
	InflationDistribution InflationDistributionParams
}

type InflationDistributionParams struct {
	StakingRewards sdk.LegacyDec
	CommunityPool  sdk.LegacyDec
}
