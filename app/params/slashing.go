package params

import (
    "time"
    "cosmossdk.io/math"
    slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

var (
    // Slashing params
    SlashingParams = slashingtypes.Params{
        SignedBlocksWindow:      100,
        MinSignedPerWindow:      math.LegacyNewDecWithPrec(500, 3), // 50%
        DowntimeJailDuration:    time.Duration(600) * time.Second,   // 10 minutes
        SlashFractionDoubleSign: math.LegacyNewDecWithPrec(50, 3),  // 5%
        SlashFractionDowntime:   math.LegacyNewDecWithPrec(10, 3),  // 1%
    }
)
