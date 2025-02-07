// app/params/token.go
package params

import (
	"math/big"

	sdk "cosmossdk.io/math"
)

const (
	TokenDenom    = "dhives"
	DisplayDenom  = "DHIVES"
	BaseDenomUnit = 18
)

var (
	// Calculate 10^18 for the denomination
	powerReduction = new(big.Int).Exp(big.NewInt(10), big.NewInt(BaseDenomUnit), nil)

	// Calculate 15M * 10^18
	TotalSupply = sdk.NewIntFromBigInt(new(big.Int).Mul(
		big.NewInt(15_000_000),
		powerReduction,
	))

	// Calculate 10M * 10^18
	CirculatingSupply = sdk.NewIntFromBigInt(new(big.Int).Mul(
		big.NewInt(10_000_000),
		powerReduction,
	))

	// Calculate 5M * 10^18
	BlockRewardsSupply = sdk.NewIntFromBigInt(new(big.Int).Mul(
		big.NewInt(5_000_000),
		powerReduction,
	))
)
