// Package chain provides the full LalaChain Cosmos SDK v0.50 application.
package chain

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/lala-protocol/lalachain/chain/x/aiadvisor"
	"github.com/lala-protocol/lalachain/chain/x/lalagov"
	"github.com/lala-protocol/lalachain/chain/x/telemetry"
)

const AppName = "lalachain"

var (
	DefaultNodeHome string

	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		sdkparams.AppModuleBasic{},
		slashing.AppModuleBasic{},
		consensus.AppModuleBasic{},
		vesting.AppModuleBasic{},
		genutil.AppModuleBasic{GenTxValidator: genutiltypes.DefaultMessageValidator},
	)

	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	}
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory instead of panicking
		DefaultNodeHome = ".lalachaind"
		return
	}
	DefaultNodeHome = filepath.Join(userHomeDir, ".lalachaind")
}

// LalaApp extends the Cosmos SDK BaseApp.
type LalaApp struct {
	*baseapp.BaseApp

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry
	keys              map[string]*storetypes.KVStoreKey

	// SDK keepers
	AccountKeeper   authkeeper.AccountKeeper
	BankKeeper      bankkeeper.Keeper
	StakingKeeper   *stakingkeeper.Keeper
	DistrKeeper     distrkeeper.Keeper
	SlashingKeeper  slashingkeeper.Keeper
	MintKeeper      mintkeeper.Keeper
	ConsensusKeeper consensuskeeper.Keeper
	ParamsKeeper    paramskeeper.Keeper

	// Custom keepers
	TelemetryKeeper telemetry.Keeper
	AIAdvisorKeeper *aiadvisor.Keeper
	LalaGovKeeper   *lalagov.Keeper

	// Epoch orchestration
	epoch epochState

	mm           *module.Manager
	configurator module.Configurator
}

// New returns a fully wired LalaChain application.
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
) *LalaApp {
	encodingConfig := MakeEncodingConfig()
	appCodec := encodingConfig.Codec
	legacyAmino := encodingConfig.Amino
	txConfig := encodingConfig.TxConfig
	interfaceRegistry := encodingConfig.InterfaceRegistry

	baseAppOptions := server.DefaultBaseappOptions(appOpts)
	bApp := baseapp.NewBaseApp(AppName, logger, db, txConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		minttypes.StoreKey,
		paramstypes.StoreKey,
		consensustypes.StoreKey,
		// Custom
		telemetry.ModuleName,
		aiadvisor.ModuleName,
		lalagov.ModuleName,
	)

	app := &LalaApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		txConfig:          txConfig,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
	}

	app.MountKVStores(keys)

	// Address codecs
	accAddrCodec := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	valAddrCodec := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())

	// ── Params ──
	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], nil)

	// ── Consensus ──
	app.ConsensusKeeper = consensuskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[consensustypes.StoreKey]),
		authtypes.NewModuleAddress("gov").String(),
		runtime.EventService{},
	)
	bApp.SetParamStore(app.ConsensusKeeper.ParamsStore)

	// ── Auth ──
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		accAddrCodec,
		sdk.Bech32MainPrefix,
		authtypes.NewModuleAddress("gov").String(),
	)

	// ── Bank ──
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper,
		BlockedModuleAccountAddrs(),
		authtypes.NewModuleAddress("gov").String(),
		logger,
	)

	// ── Staking ──
	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress("gov").String(),
		valAddrCodec,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	// ── Distribution ──
	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress("gov").String(),
	)

	// ── Slashing ──
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		app.StakingKeeper,
		authtypes.NewModuleAddress("gov").String(),
	)

	// ── Mint ──
	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[minttypes.StoreKey]),
		app.StakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress("gov").String(),
	)

	// ── Custom: Telemetry ──
	app.TelemetryKeeper = telemetry.NewKeeper(keys[telemetry.ModuleName])

	// ── Custom: AI Advisor ──
	// Advisory key secret for HMAC-SHA256 proposal signing.
	// In production, load from secure config / environment variable.
	advisorySecret := []byte("lala-testnet-advisory-secret-v1")
	app.AIAdvisorKeeper = aiadvisor.NewKeeper(
		keys[aiadvisor.ModuleName],
		"lala-ai-advisor-testnet",
		advisorySecret,
		aiadvisor.DefaultConfig(),
	)

	// ── Custom: LalaGov ──
	app.LalaGovKeeper = lalagov.NewKeeper(
		keys[lalagov.ModuleName],
		lalagov.DefaultConfig(),
		map[string][]byte{"lala-ai-advisor-testnet": advisorySecret},
	)

	// ── Module Manager ──
	app.mm = module.NewManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, nil),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, nil),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, nil),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, nil),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, nil),
		sdkparams.NewAppModule(app.ParamsKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, nil, app.interfaceRegistry),
		consensus.NewAppModule(appCodec, app.ConsensusKeeper),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app, txConfig),
		// Custom
		telemetry.NewAppModule(app.TelemetryKeeper),
		aiadvisor.NewAppModule(app.AIAdvisorKeeper),
		lalagov.NewAppModule(app.LalaGovKeeper),
	)

	app.mm.SetOrderBeginBlockers(
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		consensustypes.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		telemetry.ModuleName,
		aiadvisor.ModuleName,
		lalagov.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		stakingtypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		consensustypes.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		telemetry.ModuleName,
		aiadvisor.ModuleName,
		lalagov.ModuleName,
	)

	app.mm.SetOrderInitGenesis(
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		consensustypes.ModuleName,
		telemetry.ModuleName,
		aiadvisor.ModuleName,
		lalagov.ModuleName,
	)

	app.configurator = module.NewConfigurator(appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	if err := app.mm.RegisterServices(app.configurator); err != nil {
		panic(err)
	}

	app.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks(),
		),
	)

	anteHandler, err := NewAnteHandler(app)
	if err != nil {
		panic(fmt.Sprintf("failed to create ante handler: %v", err))
	}
	app.SetAnteHandler(anteHandler)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			panic(err)
		}
	}

	// Set a query multistore that gracefully handles IAVL version lookups.
	app.SetQueryMultiStore(newQueryableMultiStore(app.CommitMultiStore()))

	return app
}

// InitChainer initializes the chain from genesis.
func (app *LalaApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState map[string]json.RawMessage
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// BeginBlocker runs module begin-block logic.
func (app *LalaApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}

// EndBlocker runs module end-block logic.
func (app *LalaApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	res, err := app.mm.EndBlock(ctx)
	if err != nil {
		return res, err
	}
	app.runEpochEndBlocker(ctx)
	return res, nil
}

// RegisterAPIRoutes registers module REST routes with the API server.
func (app *LalaApp) RegisterAPIRoutes(apiSvr *api.Server, _ config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register CometBFT query routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node service grpc-gateway routes.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register custom REST routes for lala modules.
	app.registerCustomRESTRoutes(apiSvr.Router)
}

// RegisterTxService implements the Application interface.
func (app *LalaApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.GRPCQueryRouter(), clientCtx, app.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application interface.
func (app *LalaApp) RegisterTendermintService(clientCtx client.Context) {
	cmtApp := server.NewCometABCIWrapper(app)
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.GRPCQueryRouter(),
		app.interfaceRegistry,
		cmtApp.Query,
	)
}

// RegisterNodeService implements the Application interface.
func (app *LalaApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// ExportAppStateAndValidators exports the app state for genesis export.
func (app *LalaApp) ExportAppStateAndValidators(forZeroHeight bool, jailAllowedAddrs []string, modulesToExport []string) (servertypes.ExportedApp, error) {
	ctx := app.NewContext(true)
	genState, err := app.mm.ExportGenesisForModules(ctx, app.appCodec, modulesToExport)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	return servertypes.ExportedApp{AppState: appState}, nil
}

// BlockedModuleAccountAddrs returns addresses that should not receive funds.
func BlockedModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

func initParamsKeeper(appCodec codec.Codec, legacyAmino *codec.LegacyAmino, key *storetypes.KVStoreKey, tkey *storetypes.TransientStoreKey) paramskeeper.Keeper {
	pk := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)
	pk.Subspace(authtypes.ModuleName)
	pk.Subspace(banktypes.ModuleName)
	pk.Subspace(stakingtypes.ModuleName)
	pk.Subspace(distrtypes.ModuleName)
	pk.Subspace(slashingtypes.ModuleName)
	pk.Subspace(minttypes.ModuleName)
	return pk
}
