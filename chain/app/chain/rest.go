package chain

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	ltypes "github.com/lala-protocol/lalachain/chain/types"
)

// registerCustomRESTRoutes registers REST endpoints for custom modules
// directly on the API server's gorilla/mux Router, bypassing grpc-gateway.
func (app *LalaApp) registerCustomRESTRoutes(router *mux.Router) {
	router.HandleFunc("/lala/telemetry/v1/kpis", app.handleAllKPIs).Methods("GET")
	router.HandleFunc("/lala/lalagov/v1/history", app.handleGovHistory).Methods("GET")
	router.HandleFunc("/lala/lalagov/v1/config", app.handleGovConfig).Methods("GET")
	router.HandleFunc("/lala/aiadvisor/v1/state", app.handleAdvisorState).Methods("GET")
}

func (app *LalaApp) handleAllKPIs(w http.ResponseWriter, r *http.Request) {
	ctx, err := app.CreateQueryContext(0, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	store := ctx.KVStore(app.keys["telemetry"])
	iterator := store.Iterator([]byte("kpi/"), []byte("kpi0"))
	defer iterator.Close()

	var results []map[string]interface{}
	for ; iterator.Valid(); iterator.Next() {
		var kpis ltypes.EpochKPIs
		if err := json.Unmarshal(iterator.Value(), &kpis); err != nil {
			continue
		}
		results = append(results, map[string]interface{}{
			"epoch":               kpis.Epoch,
			"avg_block_time_ms":   kpis.AvgBlockTimeMS,
			"block_time_variance": kpis.BlockTimeVarianceMS,
			"avg_utilization":     kpis.AvgUtilization,
			"avg_base_fee":        kpis.AvgBaseFee,
			"validator_count":     kpis.ValidatorCount,
			"total_staked_ratio":  kpis.TotalStakedRatio,
			"mempool_size_trend":  kpis.MempoolSizeTrend,
			"slashing_events":     kpis.SlashingEvents,
		})
	}
	if results == nil {
		results = []map[string]interface{}{}
	}
	writeJSON(w, map[string]interface{}{"kpis": results})
}

func (app *LalaApp) handleGovHistory(w http.ResponseWriter, r *http.Request) {
	history := app.LalaGovKeeper.History()
	proposals := make([]map[string]interface{}, 0, len(history))
	for _, p := range history {
		proposals = append(proposals, map[string]interface{}{
			"proposal_id":    p.Proposal.ProposalID,
			"parameter":      p.Proposal.Parameter,
			"current_value":  p.Proposal.CurrentValue,
			"proposed_value": p.Proposal.ProposedValue,
			"rationale":      p.Proposal.Rationale,
			"votes_approve":  fmt.Sprintf("%g", p.VotesApprove),
			"votes_reject":   fmt.Sprintf("%g", p.VotesReject),
			"outcome":        p.Outcome,
		})
	}
	writeJSON(w, map[string]interface{}{"proposals": proposals})
}

func (app *LalaApp) handleGovConfig(w http.ResponseWriter, r *http.Request) {
	cfg := app.LalaGovKeeper.GetConfig()
	writeJSON(w, map[string]interface{}{
		"config": map[string]interface{}{
			"quorum":               fmt.Sprintf("%g", cfg.Quorum),
			"approval":             fmt.Sprintf("%g", cfg.Approval),
			"voting_period_epochs": cfg.VotingPeriodEpochs,
		},
	})
}

func (app *LalaApp) handleAdvisorState(w http.ResponseWriter, r *http.Request) {
	nextID, lowStreak, highStreak := app.AIAdvisorKeeper.GetState()
	writeJSON(w, map[string]interface{}{
		"state": map[string]interface{}{
			"next_proposal_id": nextID,
			"low_util_streak":  lowStreak,
			"high_util_streak": highStreak,
		},
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(v)
}
