package cosmosapp

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	protoapp "github.com/lala-protocol/lalachain/chain/app"
	ltypes "github.com/lala-protocol/lalachain/chain/types"
	"github.com/lala-protocol/lalachain/chain/x/aiadvisor"
	"github.com/lala-protocol/lalachain/chain/x/lalagov"
	"github.com/lala-protocol/lalachain/chain/x/telemetry"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

const (
	storeKeyTelemetry = "telemetry"
	storeKeyAIAdvisor = "aiadvisor"
	storeKeyGov       = "gov"
	storeKeyParams    = "params"

	chainIDLocal = "lalachain-local"
)

// Snapshot captures the persisted application-level state summary.
type Snapshot struct {
	Phase                 string               `json:"phase"`
	Validators            int                  `json:"validators"`
	Epochs                int                  `json:"epochs"`
	FinalParams           ltypes.NetworkParams `json:"final_params"`
	GovernancePassed      int                  `json:"governance_passed"`
	GovernanceFailed      int                  `json:"governance_failed"`
	KPIAverageUtilization float64              `json:"kpi_avg_utilization"`
	KPIAverageBaseFeeGwei float64              `json:"kpi_avg_base_fee_gwei"`
}

// App wires a Cosmos SDK BaseApp with ModuleManager, store-backed keepers, and a simulation prototype.
type App struct {
	baseApp *baseapp.BaseApp
	keys    map[string]*storetypes.KVStoreKey
	phase   protoapp.Phase
	mm      *module.Manager

	// Store-backed keepers owned by the module manager.
	telemetryKeeper telemetry.Keeper
	aiKeeper        *aiadvisor.Keeper
	govKeeper       *lalagov.Keeper

	// Prototype drives the deterministic simulation loop.
	prototype *protoapp.Prototype

	lastSnapshot *Snapshot
}

func New(seed int64, phase protoapp.Phase, validatorCount int) (*App, error) {
	db := dbm.NewMemDB()
	logger := log.NewNopLogger()
	txDecoder := func(txBytes []byte) (sdk.Tx, error) {
		return nil, nil
	}

	bapp := baseapp.NewBaseApp("lalachain", logger, db, txDecoder)
	keys := map[string]*storetypes.KVStoreKey{
		storeKeyTelemetry: storetypes.NewKVStoreKey(storeKeyTelemetry),
		storeKeyAIAdvisor: storetypes.NewKVStoreKey(storeKeyAIAdvisor),
		storeKeyGov:       storetypes.NewKVStoreKey(storeKeyGov),
		storeKeyParams:    storetypes.NewKVStoreKey(storeKeyParams),
	}
	bapp.MountKVStores(keys)
	if err := bapp.LoadLatestVersion(); err != nil {
		return nil, fmt.Errorf("load latest version: %w", err)
	}

	// Create store-backed keepers.
	advisoryKey := "lala-ai-advisor-phase1"
	if phase == protoapp.Phase2 {
		advisoryKey = "lala-ai-advisor-phase2"
	}

	tk := telemetry.NewKeeper(keys[storeKeyTelemetry])
	advisorySecret := []byte("lala-testnet-advisory-secret-v1")
	ak := aiadvisor.NewKeeper(keys[storeKeyAIAdvisor], advisoryKey, advisorySecret, aiadvisor.DefaultConfig())
	gk := lalagov.NewKeeper(keys[storeKeyGov], lalagov.DefaultConfig(), map[string][]byte{advisoryKey: advisorySecret})

	// Create modules.
	telemetryMod := telemetry.NewAppModule(tk)
	aiMod := aiadvisor.NewAppModule(ak)
	govMod := lalagov.NewAppModule(gk)

	// Create ModuleManager with ordering: telemetry → aiadvisor → gov.
	mm := module.NewManager(telemetryMod, aiMod, govMod)
	mm.SetOrderBeginBlockers(telemetry.ModuleName, aiadvisor.ModuleName, lalagov.ModuleName)
	mm.SetOrderEndBlockers(telemetry.ModuleName, aiadvisor.ModuleName, lalagov.ModuleName)

	// Create prototype with nil store keys (in-memory simulation driver).
	prototype, err := newPrototype(seed, phase, validatorCount)
	if err != nil {
		return nil, err
	}

	app := &App{
		baseApp:         bapp,
		keys:            keys,
		phase:           phase,
		mm:              mm,
		telemetryKeeper: tk,
		aiKeeper:        ak,
		govKeeper:       gk,
		prototype:       prototype,
	}

	// Wire epoch hook so module lifecycle runs at each epoch boundary.
	prototype.SetEpochHook(app)

	return app, nil
}

func newPrototype(seed int64, phase protoapp.Phase, validatorCount int) (*protoapp.Prototype, error) {
	switch phase {
	case protoapp.Phase1:
		return protoapp.NewPrototype(seed), nil
	case protoapp.Phase2:
		return protoapp.NewPhase2Prototype(seed, validatorCount), nil
	default:
		return nil, fmt.Errorf("unsupported phase %q", phase)
	}
}

// BeforeEpoch implements protoapp.EpochHook — called before each epoch starts.
// Triggers module BeginBlockers (gov loads state from store).
func (a *App) BeforeEpoch(epoch int64) {
	ctx := a.newContext(epoch)
	_, _ = a.mm.BeginBlock(ctx)
}

// AfterEpoch implements protoapp.EpochHook — called after each epoch completes.
// Persists KPIs and triggers module EndBlockers (all modules save state).
func (a *App) AfterEpoch(epoch int64) {
	ctx := a.newContext(epoch)

	// Persist latest KPIs from prototype into the telemetry store.
	kpiHistory := a.prototype.KPIHistory()
	if len(kpiHistory) > 0 {
		latestKPI := kpiHistory[len(kpiHistory)-1]
		a.telemetryKeeper.PersistKPIs(ctx, latestKPI)
	}

	// Sync aiadvisor state from prototype's keeper into store-backed keeper.
	// The prototype's AI keeper has the authoritative in-memory state.
	a.aiKeeper.SaveState(ctx)

	// Trigger module EndBlockers.
	_, _ = a.mm.EndBlock(ctx)
}

func (a *App) Run(epochs int, out io.Writer) error {
	if a == nil || a.baseApp == nil || a.prototype == nil {
		return errors.New("cosmos app is not initialized")
	}
	if epochs <= 0 {
		return errors.New("epochs must be > 0")
	}

	ctx := a.newContext(1)
	a.persistBootMetadata(ctx, epochs)

	if err := a.prototype.Run(epochs, out); err != nil {
		return err
	}

	snapshot := a.buildSnapshot(epochs)
	a.lastSnapshot = &snapshot
	a.persistSnapshot(ctx, snapshot)

	return nil
}

func (a *App) Snapshot() (Snapshot, bool) {
	if a == nil || a.lastSnapshot == nil {
		return Snapshot{}, false
	}
	return *a.lastSnapshot, true
}

func (a *App) StateHash() string {
	snapshot, ok := a.Snapshot()
	if !ok {
		return ""
	}
	payload, err := json.Marshal(snapshot)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

// ModuleManager returns the app's module manager for external inspection/testing.
func (a *App) ModuleManager() *module.Manager {
	return a.mm
}

func (a *App) newContext(height int64) sdk.Context {
	header := cmtproto.Header{
		ChainID: chainIDLocal,
		Height:  height,
		Time:    time.Now().UTC(),
	}
	return a.baseApp.NewUncachedContext(false, header)
}

func (a *App) persistBootMetadata(ctx sdk.Context, epochs int) {
	paramsStore := ctx.KVStore(a.keys[storeKeyParams])
	paramsStore.Set([]byte("phase"), []byte(a.phase))
	paramsStore.Set([]byte("epochs"), []byte(fmt.Sprintf("%d", epochs)))
}

func (a *App) buildSnapshot(epochs int) Snapshot {
	history := a.prototype.History()
	passed := 0
	failed := 0
	for _, item := range history {
		if item.Outcome == "passed" {
			passed++
		} else {
			failed++
		}
	}

	kpiHistory := a.prototype.KPIHistory()
	utilSum := 0.0
	feeSum := 0.0
	for _, kpi := range kpiHistory {
		utilSum += kpi.AvgUtilization
		feeSum += float64(kpi.AvgBaseFee) / 1e9
	}
	avgUtil := 0.0
	avgFee := 0.0
	if len(kpiHistory) > 0 {
		avgUtil = utilSum / float64(len(kpiHistory))
		avgFee = feeSum / float64(len(kpiHistory))
	}

	return Snapshot{
		Phase:                 string(a.phase),
		Validators:            validatorCount(kpiHistory),
		Epochs:                epochs,
		FinalParams:           a.prototype.Params(),
		GovernancePassed:      passed,
		GovernanceFailed:      failed,
		KPIAverageUtilization: avgUtil,
		KPIAverageBaseFeeGwei: avgFee,
	}
}

func validatorCount(kpis []ltypes.EpochKPIs) int {
	if len(kpis) == 0 {
		return 0
	}
	return kpis[len(kpis)-1].ValidatorCount
}

func (a *App) persistSnapshot(ctx sdk.Context, snapshot Snapshot) {
	payload, err := json.Marshal(snapshot)
	if err != nil {
		return
	}

	telemetryStore := ctx.KVStore(a.keys[storeKeyTelemetry])
	telemetryStore.Set([]byte("kpi_avg_utilization"), []byte(fmt.Sprintf("%.6f", snapshot.KPIAverageUtilization)))
	telemetryStore.Set([]byte("kpi_avg_base_fee_gwei"), []byte(fmt.Sprintf("%.6f", snapshot.KPIAverageBaseFeeGwei)))

	aiStore := ctx.KVStore(a.keys[storeKeyAIAdvisor])
	aiStore.Set([]byte("phase"), []byte(snapshot.Phase))
	aiStore.Set([]byte("snapshot"), payload)

	govStore := ctx.KVStore(a.keys[storeKeyGov])
	govStore.Set([]byte("passed"), []byte(fmt.Sprintf("%d", snapshot.GovernancePassed)))
	govStore.Set([]byte("failed"), []byte(fmt.Sprintf("%d", snapshot.GovernanceFailed)))

	paramsStore := ctx.KVStore(a.keys[storeKeyParams])
	paramsStore.Set([]byte("snapshot"), payload)
	paramsStore.Set([]byte("state_hash"), []byte(a.StateHash()))
}
