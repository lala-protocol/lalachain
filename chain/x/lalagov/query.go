package lalagov

import (
	"context"
	"fmt"

	aiadvisorv1 "github.com/lala-protocol/lalachain/chain/api/lala/aiadvisor/v1"
	lalagovv1 "github.com/lala-protocol/lalachain/chain/api/lala/lalagov/v1"
)

// queryServer implements the lalagov gRPC QueryServer.
type queryServer struct {
	lalagovv1.UnimplementedQueryServer
	keeper *Keeper
}

func NewQueryServer(keeper *Keeper) lalagovv1.QueryServer {
	return &queryServer{keeper: keeper}
}

func (s *queryServer) Config(_ context.Context, _ *lalagovv1.QueryConfigRequest) (*lalagovv1.QueryConfigResponse, error) {
	s.keeper.mu.RLock()
	cfg := s.keeper.config
	s.keeper.mu.RUnlock()

	return &lalagovv1.QueryConfigResponse{
		Config: &lalagovv1.GovConfig{
			Quorum:             fmt.Sprintf("%g", cfg.Quorum),
			Approval:           fmt.Sprintf("%g", cfg.Approval),
			VotingPeriodEpochs: cfg.VotingPeriodEpochs,
		},
	}, nil
}

func (s *queryServer) History(_ context.Context, _ *lalagovv1.QueryHistoryRequest) (*lalagovv1.QueryHistoryResponse, error) {
	history := s.keeper.History()

	proposals := make([]*lalagovv1.ResolvedProposal, 0, len(history))
	for _, rp := range history {
		proposals = append(proposals, &lalagovv1.ResolvedProposal{
			Proposal: &aiadvisorv1.AIProposal{
				ProposalId:        rp.Proposal.ProposalID,
				EpochSubmitted:    rp.Proposal.EpochSubmitted,
				Parameter:         rp.Proposal.Parameter,
				CurrentValue:      rp.Proposal.CurrentValue,
				ProposedValue:     rp.Proposal.ProposedValue,
				Rationale:         rp.Proposal.Rationale,
				ActivationEpoch:   rp.Proposal.ActivationEpoch,
				AdvisoryKey:       rp.Proposal.AdvisoryKey,
				AdvisorySignature: rp.Proposal.AdvisorySignature,
			},
			VotesApprove: fmt.Sprintf("%g", rp.VotesApprove),
			VotesReject:  fmt.Sprintf("%g", rp.VotesReject),
			Outcome:      rp.Outcome,
		})
	}

	return &lalagovv1.QueryHistoryResponse{Proposals: proposals}, nil
}
