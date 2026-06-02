package aiadvisor

import (
	"context"

	aiadvisorv1 "github.com/lala-protocol/lalachain/chain/api/lala/aiadvisor/v1"
)

// queryServer implements the aiadvisor gRPC QueryServer.
type queryServer struct {
	aiadvisorv1.UnimplementedQueryServer
	keeper *Keeper
}

func NewQueryServer(keeper *Keeper) aiadvisorv1.QueryServer {
	return &queryServer{keeper: keeper}
}

func (s *queryServer) State(_ context.Context, _ *aiadvisorv1.QueryStateRequest) (*aiadvisorv1.QueryStateResponse, error) {
	s.keeper.mu.Lock()
	state := &aiadvisorv1.AdvisorState{
		NextProposalId: s.keeper.nextProposalID,
		LowUtilStreak:  s.keeper.lowUtilStreak,
		HighUtilStreak: s.keeper.highUtilStreak,
	}
	s.keeper.mu.Unlock()

	return &aiadvisorv1.QueryStateResponse{State: state}, nil
}
