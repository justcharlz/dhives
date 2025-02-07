package params

import (
    "cosmossdk.io/math"
)

var (
    // MainnetMinGasPrices defines 0.1 dhives as minimum gas price
    MainnetMinGasPrices = math.LegacyNewDec(100_000_000_000) // 0.1 dhives
    MainnetMinGasMultiplier = math.LegacyNewDecWithPrec(5, 1) // 0.5
)
