package main

import (
    "os"

    svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/justcharlz/dhives/app"
    cmdcfg "github.com/justcharlz/dhives/cmd/config"
)

func main() {
    cfg := sdk.GetConfig()
    cmdcfg.SetBech32Prefixes(cfg)
    cmdcfg.SetBip44CoinType(cfg)
    cfg.Seal()

    rootCmd, _ := NewRootCmd()
    if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome, app.DefaultNodeHome); err != nil {
        os.Exit(1)
    }
}
