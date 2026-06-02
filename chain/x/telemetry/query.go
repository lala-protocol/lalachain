package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ltypes "github.com/lala-protocol/lalachain/chain/types"

	telemetryv1 "github.com/lala-protocol/lalachain/chain/api/lala/telemetry/v1"
)

// queryServer implements the telemetry gRPC QueryServer.
type queryServer struct {
	telemetryv1.UnimplementedQueryServer
	keeper Keeper
}

func NewQueryServer(keeper Keeper) telemetryv1.QueryServer {
	return &queryServer{keeper: keeper}
}

func (s *queryServer) KPIs(ctx context.Context, req *telemetryv1.QueryKPIsRequest) (*telemetryv1.QueryKPIsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kpis, found := s.keeper.GetKPIs(sdkCtx, req.Epoch)
	if !found {
		return &telemetryv1.QueryKPIsResponse{}, nil
	}
	return &telemetryv1.QueryKPIsResponse{
		Kpis: convertEpochKPIs(kpis),
	}, nil
}

func (s *queryServer) AllKPIs(ctx context.Context, _ *telemetryv1.QueryAllKPIsRequest) (*telemetryv1.QueryAllKPIsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if s.keeper.storeKey == nil {
		return &telemetryv1.QueryAllKPIsResponse{}, nil
	}
	store := sdkCtx.KVStore(s.keeper.storeKey)
	iterator := store.Iterator([]byte("kpi/"), []byte("kpi0")) // "kpi0" is just after "kpi/" lexicographically
	defer iterator.Close()

	var results []*telemetryv1.EpochKPIs
	for ; iterator.Valid(); iterator.Next() {
		var kpis ltypes.EpochKPIs
		if err := json.Unmarshal(iterator.Value(), &kpis); err != nil {
			continue
		}
		results = append(results, convertEpochKPIs(kpis))
	}
	return &telemetryv1.QueryAllKPIsResponse{Kpis: results}, nil
}

func convertEpochKPIs(k ltypes.EpochKPIs) *telemetryv1.EpochKPIs {
	return &telemetryv1.EpochKPIs{
		Epoch:             k.Epoch,
		AvgBlockTimeMs:    formatFloat(k.AvgBlockTimeMS),
		BlockTimeVariance: formatFloat(k.BlockTimeVarianceMS),
		AvgUtilization:    formatFloat(k.AvgUtilization),
		AvgBaseFee:        k.AvgBaseFee,
		ValidatorCount:    int32(k.ValidatorCount),
		TotalStakedRatio:  formatFloat(k.TotalStakedRatio),
		MempoolSizeTrend:  formatFloat(k.MempoolSizeTrend),
		SlashingEvents:    k.SlashingEvents,
	}
}

func formatFloat(f float64) string {
	s := fmt.Sprintf("%g", f)
	if !strings.Contains(s, ".") && !strings.Contains(s, "e") {
		s += ".0"
	}
	return s
}
