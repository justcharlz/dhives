package params

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

var (
	// Governance params
	GovParams = govv1.Params{
		MinDeposit:    sdk.NewCoins(sdk.NewCoin(TokenDenom, math.NewInt(10_000_000_000_000_000_000))), // 10 DHIVES
		VotingPeriod:  time.Duration(172800) * time.Second,                                            // 2 days
		Quorum:        math.LegacyNewDecWithPrec(334, 3),                                              // 33.4%
		Threshold:     math.LegacyNewDecWithPrec(500, 3),                                              // 50%
		VetoThreshold: math.LegacyNewDecWithPrec(334, 3),                                              // 33.4%
	}
)
