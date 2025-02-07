package params

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

var (
	// Define voting period separately
	votingPeriod = time.Duration(172800 * time.Second) // 2 days

	// Governance params
	GovParams = govv1.Params{
		MinDeposit:                 sdk.NewCoins(sdk.NewCoin(TokenDenom, math.NewIntFromUint64(10000000000000000000))), // 10 DHIVES
		VotingPeriod:               &votingPeriod,                                                                      // 2 days
		Quorum:                     "0.334",                                                                            // 33.4%
		Threshold:                  "0.500",                                                                            // 50%
		VetoThreshold:              "0.334",                                                                            // 33.4%
		MinInitialDepositRatio:     "0.000000000000000000",                                                             // 0%
		ProposalCancelRatio:        "0.500",                                                                            // 50%
		ProposalCancelDest:         "",
		BurnVoteQuorum:             false,
		BurnProposalDepositPrevote: false,
		BurnVoteVeto:               true,
	}
)
