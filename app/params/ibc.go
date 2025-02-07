package params

import (
    ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

var (
    // IBC Transfer params
    IBCTransferParams = ibctransfertypes.Params{
        SendEnabled:    true,
        ReceiveEnabled: true,
    }
)
