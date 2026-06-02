package chain

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ltypes "github.com/lala-protocol/lalachain/chain/types"
)

const EpochLength = 10

// epochState tracks in-memory block data across an epoch for KPI computation.
type epochState struct {
	blocks        []ltypes.FinalizedBlock
	lastBlockTime time.Time
	currentEpoch  int64
}

// runEpochEndBlocker is called from EndBlocker on every block.
// It accumulates block telemetry and triggers epoch-boundary logic.
func (app *LalaApp) runEpochEndBlocker(ctx sdk.Context) {
	height := ctx.BlockHeight()
	if height <= 1 {
		// Genesis block – just record the time.
		app.epoch.lastBlockTime = ctx.BlockTime()
		return
	}

	// Compute block time delta from previous block.
	blockTimeMS := int64(0)
	if !app.epoch.lastBlockTime.IsZero() {
		blockTimeMS = ctx.BlockTime().Sub(app.epoch.lastBlockTime).Milliseconds()
	}
	app.epoch.lastBlockTime = ctx.BlockTime()

	// On a testnet with no real transactions, gas used is 0.
	// We use the deliver-tx gas from the context (actual consumed gas this block).
	gasUsed := ctx.BlockGasMeter().GasConsumed()

	// Simulate EIP-1559 base fee decay: when utilization is 0, the fee drops.
	// Start from the param value and decay 12.5% per empty block in this epoch.
	baseFee := ltypes.DefaultNetworkParams().BaseFeePerGas
	if gasUsed == 0 {
		// Each consecutive empty block decays base fee by 12.5%
		// After ~6 empty blocks, fee drops below MinFeeTarget (0.8 Gwei).
		decayFactor := int64(1)
		for i := 0; i < len(app.epoch.blocks); i++ {
			decayFactor++
		}
		baseFee = baseFee * 7 / (7 + decayFactor) // progressive decay
	}

	app.epoch.blocks = append(app.epoch.blocks, ltypes.FinalizedBlock{
		BlockTimeMS: blockTimeMS,
		GasUsed:     int64(gasUsed),
		BaseFee:     baseFee,
	})

	// Check epoch boundary.
	if height%EpochLength != 0 {
		return
	}

	epoch := height / EpochLength
	app.epoch.currentEpoch = epoch

	ctx.Logger().Info("epoch boundary reached", "epoch", epoch, "height", height)

	// ── 1. Compute KPIs ──
	params := ltypes.DefaultNetworkParams()
	validators := []ltypes.Validator{{Address: "testnode", Stake: 1.0}}

	kpis, err := app.TelemetryKeeper.CalculateEpochKPIs(
		epoch,
		app.epoch.blocks,
		params,
		validators,
		0.0, // mempool trend
		0,   // slashing events
	)
	if err != nil {
		ctx.Logger().Error("epoch KPI computation failed", "error", err)
		app.epoch.blocks = nil
		return
	}

	// Persist KPIs to the store.
	app.TelemetryKeeper.PersistKPIs(ctx, kpis)
	ctx.Logger().Info("epoch KPIs persisted",
		"epoch", epoch,
		"avg_util", kpis.AvgUtilization,
		"avg_block_time_ms", kpis.AvgBlockTimeMS,
	)

	// ── 2. AI Advisor: analyze and maybe propose ──
	proposal := app.AIAdvisorKeeper.AnalyzeAndPropose(kpis, params)
	if proposal != nil {
		ctx.Logger().Info("AI advisor generated proposal",
			"id", proposal.ProposalID,
			"param", proposal.Parameter,
			"current", proposal.CurrentValue,
			"proposed", proposal.ProposedValue,
			"rationale", proposal.Rationale,
		)

		// Register the proposal in governance.
		if err := app.LalaGovKeeper.RegisterAIProposal(*proposal); err != nil {
			ctx.Logger().Error("failed to register AI proposal", "error", err)
		} else {
			// Auto-vote on single-validator testnet (the only validator approves).
			_ = app.LalaGovKeeper.CastVote(
				proposal.ProposalID,
				ltypes.Validator{Address: "testnode", Stake: 1.0},
				true,
			)
		}
	}

	// ── 3. Governance: resolve proposals whose voting period ended ──
	resolved := app.LalaGovKeeper.EndEpoch(epoch, 1.0)
	for _, r := range resolved {
		ctx.Logger().Info("proposal resolved",
			"id", r.Proposal.ProposalID,
			"outcome", r.Outcome,
		)
	}

	// Persist keeper states after epoch mutations (module EndBlockers run before this).
	app.AIAdvisorKeeper.SaveState(ctx)
	app.LalaGovKeeper.SaveState(ctx)

	// Reset block buffer for next epoch.
	app.epoch.blocks = nil
}
