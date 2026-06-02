package telemetry

import (
	"context"
	"encoding/json"

	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	telemetryv1 "github.com/lala-protocol/lalachain/chain/api/lala/telemetry/v1"
)

const ModuleName = "telemetry"

var (
	_ module.AppModule        = AppModule{}
	_ appmodule.AppModule     = AppModule{}
	_ appmodule.HasEndBlocker = AppModule{}
)

// AppModule implements the Cosmos SDK module interface for telemetry.
type AppModule struct {
	keeper Keeper
}

func NewAppModule(keeper Keeper) AppModule {
	return AppModule{keeper: keeper}
}

func (AppModule) Name() string                                                    { return ModuleName }
func (AppModule) IsOnePerModuleType()                                             {}
func (AppModule) IsAppModule()                                                    {}
func (AppModule) ConsensusVersion() uint64                                        { return 1 }
func (AppModule) RegisterLegacyAminoCodec(*codec.LegacyAmino)                     {}
func (AppModule) RegisterInterfaces(codectypes.InterfaceRegistry)                 {}
func (AppModule) RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux)     {}
func (AppModule) DefaultGenesis(codec.JSONCodec) json.RawMessage                  { return nil }
func (AppModule) ValidateGenesis(codec.JSONCodec, client.TxEncodingConfig, json.RawMessage) error {
	return nil
}
func (am AppModule) RegisterServices(cfg module.Configurator) {
	telemetryv1.RegisterQueryServer(cfg.QueryServer(), NewQueryServer(am.keeper))
}

// EndBlock persists KPI state. The actual KPI computation is driven by the
// app-level epoch loop which calls PersistKPIs directly.
func (am AppModule) EndBlock(ctx context.Context) error {
	_ = sdk.UnwrapSDKContext(ctx)
	return nil
}
