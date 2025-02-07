package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	"cosmossdk.io/log"
	tmcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/justcharlz/dhives/app"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func NewRootCmd() (*cobra.Command, error) {
	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode("block").
		WithHomeDir(app.DefaultNodeHome).
		WithViper("") // Initialize viper

	rootCmd := &cobra.Command{
		Use:   "dhivesd",
		Short: "Dhives Daemon",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// Set the client context
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd, "", nil)
		},
	}

	// Add subcommands
	rootCmd.AddCommand(
		genutilcli.InitCmd(module.NewBasicManager(), app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome, genutiltypes.DefaultMessageValidator, encodingConfig.TxConfig.SigningContext().ValidatorAddressCodec()),
		genutilcli.GenTxCmd(
			module.NewBasicManager(),
			encodingConfig.TxConfig,
			banktypes.GenesisBalancesIterator{},
			app.DefaultNodeHome,
			encodingConfig.TxConfig.SigningContext().ValidatorAddressCodec(),
		),
		genutilcli.ValidateGenesisCmd(module.NewBasicManager()),
		AddGenesisAccountCmd(app.DefaultNodeHome), // Add genesis account command
		keys.Commands(),
		server.StatusCommand(),
	)

	// Add flags
	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, createAppAndExport, addModuleInitFlags)

	return rootCmd, nil
}

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisAccountCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add a genesis account to genesis.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			config := sdk.GetConfig()
			config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			// Create and save the genesis account
			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return err
			}

			authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)
			accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return err
			}

			if accs.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}

			// Add the new account to the set of genesis accounts
			baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)
			accs = append(accs, baseAccount)

			genAccs, err := authtypes.PackAccounts(accs)
			if err != nil {
				return err
			}

			authGenState.Accounts = genAccs
			appState[authtypes.ModuleName] = cdc.MustMarshalJSON(&authGenState)

			bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
			bankGenState.Balances = append(bankGenState.Balances, banktypes.Balance{
				Address: addr.String(),
				Coins:   coins,
			})
			appState[banktypes.ModuleName] = cdc.MustMarshalJSON(bankGenState)

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return err
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "node's home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)

	return app.NewEvmos(
		logger,
		db,
		traceStore,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		appOpts,
		nil,
		baseappOptions...,
	)
}

func createAppAndExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	var dhivesApp *app.Evmos
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))

	if height != -1 {
		dhivesApp = app.NewEvmos(
			logger,
			db,
			traceStore,
			false,
			map[int64]bool{},
			homePath,
			uint(1),
			appOpts,
			nil,
			server.DefaultBaseappOptions(appOpts)...,
		)
		if err := dhivesApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		dhivesApp = app.NewEvmos(
			logger,
			db,
			traceStore,
			true,
			map[int64]bool{},
			homePath,
			uint(1),
			appOpts,
			nil,
			server.DefaultBaseappOptions(appOpts)...,
		)
	}
	return dhivesApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

func initAppConfig() (string, *serverconfig.Config) {
	srvCfg := serverconfig.DefaultConfig()
	srvCfg.MinGasPrices = "0.1dhives"

	return serverconfig.DefaultConfigTemplate, srvCfg
}

func initTendermintConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()
	return cfg
}
