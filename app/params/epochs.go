package params

import (
    "time"
    epochstypes "github.com/justcharlz/dhives/x/epochs/types"
)

var (
    // Epochs configuration
    EpochsConfig = []epochstypes.EpochInfo{
        {
            Identifier:              "day",
            StartTime:              time.Time{},
            Duration:               time.Duration(24) * time.Hour,
            CurrentEpoch:           0,
            CurrentEpochStartTime:  time.Time{},
            EpochCountingStarted:   false,
            CurrentEpochStartHeight: 0,
        },
        {
            Identifier:              "week",
            StartTime:              time.Time{},
            Duration:               time.Duration(7) * 24 * time.Hour,
            CurrentEpoch:           0,
            CurrentEpochStartTime:  time.Time{},
            EpochCountingStarted:   false,
            CurrentEpochStartHeight: 0,
        },
    }
)
