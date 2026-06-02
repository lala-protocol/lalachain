package cmd

import (
	"io"
	"os"
	"time"

	"cosmossdk.io/log"

	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"

	"github.com/lala-protocol/lalachain/chain/app/chain"
)

// NewRootCmd creates the root command for lalachaind.
func NewRootCmd() *cobra.Command {
	encodingConfig := chain.MakeEncodingConfig()

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(chain.DefaultNodeHome).
		WithViper("LALA")

	rootCmd := &cobra.Command{
		Use:   "lalachaind",
		Short: "LalaChain - AI-governed blockchain daemon",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customCMTConfig := initCometBFTConfig()
			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customCMTConfig)
		},
	}

	initRootCmd(rootCmd, chain.ModuleBasics)
	return rootCmd
}

func initRootCmd(rootCmd *cobra.Command, mb module.BasicManager) {
	valAddrCodec := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	accAddrCodec := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	rootCmd.AddCommand(
		genutilcli.InitCmd(mb, chain.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(
			banktypes.GenesisBalancesIterator{},
			chain.DefaultNodeHome,
			genutiltypes.DefaultMessageValidator,
			valAddrCodec,
		),
		genutilcli.GenTxCmd(
			mb,
			chain.MakeEncodingConfig().TxConfig,
			banktypes.GenesisBalancesIterator{},
			chain.DefaultNodeHome,
			valAddrCodec,
		),
		genutilcli.ValidateGenesisCmd(mb),
		genutilcli.AddGenesisAccountCmd(chain.DefaultNodeHome, accAddrCodec),
		debug.Cmd(),
		keys.Commands(),
	)

	// Query commands
	queryCmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(
		rpc.QueryEventForTxCmd(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)
	rootCmd.AddCommand(queryCmd)

	// Transaction commands
	txCmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetBroadcastCommand(),
	)
	rootCmd.AddCommand(txCmd)

	server.AddCommands(rootCmd, chain.DefaultNodeHome, newApp, exportApp, addModuleInitFlags)
}

func addModuleInitFlags(startCmd *cobra.Command) {}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	return chain.New(logger, db, traceStore, true, appOpts)
}

func exportApp(logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string, appOpts servertypes.AppOptions, modulesToExport []string) (servertypes.ExportedApp, error) {
	app := chain.New(logger, db, traceStore, false, appOpts)
	return app.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}

func initAppConfig() (string, interface{}) {
	srvCfg := serverconfig.DefaultConfig()
	srvCfg.MinGasPrices = "0ulala"
	return serverconfig.DefaultConfigTemplate, srvCfg
}

func initCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()
	cfg.Consensus.TimeoutCommit = 5 * time.Second
	return cfg
}
