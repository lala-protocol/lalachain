package chain

import (
	"encoding/json"
	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	// BondDenom is the staking denomination for $LALA.
	BondDenom = "ulala"

	// DisplayDenom is the human-readable denomination.
	DisplayDenom = "LALA"

	// InitialSupply is 100M LALA = 100_000_000 * 10^6 ulala.
	InitialSupply = 100_000_000_000_000
)

// DefaultGenesisState returns the default genesis with $LALA tokenomics.
func DefaultGenesisState(cdc codec.Codec) map[string]json.RawMessage {
	genesis := ModuleBasics.DefaultGenesis(cdc)

	// ── Staking ──
	var stakingGenesis stakingtypes.GenesisState
	cdc.MustUnmarshalJSON(genesis[stakingtypes.ModuleName], &stakingGenesis)
	stakingGenesis.Params.BondDenom = BondDenom
	stakingGenesis.Params.UnbondingTime = 21 * 24 * time.Hour
	stakingGenesis.Params.MaxValidators = 100
	stakingGenesis.Params.MaxEntries = 7
	stakingGenesis.Params.HistoricalEntries = 10000
	stakingGenesis.Params.MinCommissionRate = sdkmath.LegacyNewDecWithPrec(5, 2)
	genesis[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingGenesis)

	// ── Mint: 13% initial, 7-20% range, 67% goal bonded ──
	var mintGenesis minttypes.GenesisState
	cdc.MustUnmarshalJSON(genesis[minttypes.ModuleName], &mintGenesis)
	mintGenesis.Minter.Inflation = sdkmath.LegacyNewDecWithPrec(13, 2)
	mintGenesis.Minter.AnnualProvisions = sdkmath.LegacyNewDec(0)
	mintGenesis.Params.MintDenom = BondDenom
	mintGenesis.Params.InflationRateChange = sdkmath.LegacyNewDecWithPrec(13, 2)
	mintGenesis.Params.InflationMax = sdkmath.LegacyNewDecWithPrec(20, 2)
	mintGenesis.Params.InflationMin = sdkmath.LegacyNewDecWithPrec(7, 2)
	mintGenesis.Params.GoalBonded = sdkmath.LegacyNewDecWithPrec(67, 2)
	mintGenesis.Params.BlocksPerYear = 6_311_520
	genesis[minttypes.ModuleName] = cdc.MustMarshalJSON(&mintGenesis)

	// ── Distribution ──
	var distrGenesis distrtypes.GenesisState
	cdc.MustUnmarshalJSON(genesis[distrtypes.ModuleName], &distrGenesis)
	distrGenesis.Params.CommunityTax = sdkmath.LegacyNewDecWithPrec(2, 2)
	distrGenesis.Params.WithdrawAddrEnabled = true
	genesis[distrtypes.ModuleName] = cdc.MustMarshalJSON(&distrGenesis)

	// ── Slashing ──
	var slashGenesis slashingtypes.GenesisState
	cdc.MustUnmarshalJSON(genesis[slashingtypes.ModuleName], &slashGenesis)
	slashGenesis.Params.SignedBlocksWindow = 100
	slashGenesis.Params.MinSignedPerWindow = sdkmath.LegacyNewDecWithPrec(5, 1)
	slashGenesis.Params.DowntimeJailDuration = 10 * time.Minute
	slashGenesis.Params.SlashFractionDoubleSign = sdkmath.LegacyNewDecWithPrec(5, 2)
	slashGenesis.Params.SlashFractionDowntime = sdkmath.LegacyNewDecWithPrec(1, 4)
	genesis[slashingtypes.ModuleName] = cdc.MustMarshalJSON(&slashGenesis)

	// ── Bank: $LALA denom metadata ──
	var bankGenesis banktypes.GenesisState
	cdc.MustUnmarshalJSON(genesis[banktypes.ModuleName], &bankGenesis)
	bankGenesis.DenomMetadata = []banktypes.Metadata{
		{
			Description: "The native staking and governance token of LalaChain",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: BondDenom, Exponent: 0, Aliases: []string{"microlala"}},
				{Denom: "mlala", Exponent: 3, Aliases: []string{"millilala"}},
				{Denom: DisplayDenom, Exponent: 6, Aliases: nil},
			},
			Base:    BondDenom,
			Display: DisplayDenom,
			Name:    "LALA",
			Symbol:  "LALA",
		},
	}
	genesis[banktypes.ModuleName] = cdc.MustMarshalJSON(&bankGenesis)

	return genesis
}

// FaucetCoins returns the amount given by the faucet per request (10 LALA).
func FaucetCoins() sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(BondDenom, sdkmath.NewInt(10_000_000)))
}
