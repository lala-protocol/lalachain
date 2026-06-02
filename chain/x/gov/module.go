package gov

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
)

const ModuleName = "gov"

var (
	_ module.AppModule          = AppModule{}
	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
	_ appmodule.HasEndBlocker   = AppModule{}
)

// AppModule implements the Cosmos SDK module interface for governance.
type AppModule struct {
	keeper *Keeper
}

func NewAppModule(keeper *Keeper) AppModule {
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
func (AppModule) RegisterServices(module.Configurator) {}

// BeginBlock loads governance state from the store at the start of each block.
func (am AppModule) BeginBlock(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	am.keeper.LoadState(sdkCtx)
	return nil
}

// EndBlock persists governance state to the store at the end of each block.
func (am AppModule) EndBlock(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	am.keeper.SaveState(sdkCtx)
	return nil
}
