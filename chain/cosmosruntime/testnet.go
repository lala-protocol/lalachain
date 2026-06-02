package cosmosruntime

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lala-protocol/lalachain/chain/cosmosapp"
)

type LocalTestnetConfig struct {
	Phase      string
	Epochs     int
	Seed       int64
	Validators int
	Nodes      int
}

type NodeReport struct {
	NodeID    string
	StateHash string
	Snapshot  cosmosapp.Snapshot
	Output    string
}

type LocalTestnetReport struct {
	Config        LocalTestnetConfig
	Nodes         []NodeReport
	Consensus     bool
	Divergences   []string
	ReferenceHash string
}

func RunLocalTestnet(cfg LocalTestnetConfig, out io.Writer) (LocalTestnetReport, error) {
	if cfg.Epochs <= 0 {
		return LocalTestnetReport{}, errors.New("epochs must be > 0")
	}
	if cfg.Validators < 4 || cfg.Validators > 10 {
		return LocalTestnetReport{}, fmt.Errorf("validators must be in [4,10], got %d", cfg.Validators)
	}
	nodeCount := cfg.Nodes
	if nodeCount == 0 {
		nodeCount = cfg.Validators
	}
	if nodeCount <= 0 {
		return LocalTestnetReport{}, fmt.Errorf("nodes must be > 0, got %d", nodeCount)
	}

	report := LocalTestnetReport{
		Config: cfg,
		Nodes:  make([]NodeReport, 0, nodeCount),
	}

	for i := 0; i < nodeCount; i++ {
		runner, err := NewRunner(cfg.Seed, cfg.Phase, cfg.Validators)
		if err != nil {
			return LocalTestnetReport{}, fmt.Errorf("node %d init failed: %w", i+1, err)
		}

		var nodeOut bytes.Buffer
		if err := runner.Run(sdk.Context{}, cfg.Epochs, &nodeOut); err != nil {
			return LocalTestnetReport{}, fmt.Errorf("node %d run failed: %w", i+1, err)
		}

		snapshot, ok := runner.Snapshot()
		if !ok {
			return LocalTestnetReport{}, fmt.Errorf("node %d missing snapshot", i+1)
		}

		node := NodeReport{
			NodeID:    fmt.Sprintf("node-%02d", i+1),
			StateHash: runner.StateHash(),
			Snapshot:  snapshot,
			Output:    nodeOut.String(),
		}
		report.Nodes = append(report.Nodes, node)

		if i == 0 {
			report.ReferenceHash = node.StateHash
			continue
		}
		if node.StateHash != report.ReferenceHash {
			report.Divergences = append(
				report.Divergences,
				fmt.Sprintf("%s hash %s != %s", node.NodeID, node.StateHash, report.ReferenceHash),
			)
		}
	}

	report.Consensus = len(report.Divergences) == 0
	printLocalTestnetReport(report, out)

	return report, nil
}

func printLocalTestnetReport(report LocalTestnetReport, out io.Writer) {
	fmt.Fprintln(out, "======================================================================")
	fmt.Fprintln(out, " LOCAL MULTI-VALIDATOR COSMOS TESTNET")
	fmt.Fprintln(out, "======================================================================")
	fmt.Fprintf(
		out,
		"phase=%s validators=%d nodes=%d epochs=%d seed=%d\n",
		report.Config.Phase,
		report.Config.Validators,
		len(report.Nodes),
		report.Config.Epochs,
		report.Config.Seed,
	)

	for _, node := range report.Nodes {
		fmt.Fprintf(
			out,
			"[%s] hash=%s passed=%d failed=%d util=%.1f%% fee=%.3f gwei\n",
			node.NodeID,
			shortHash(node.StateHash),
			node.Snapshot.GovernancePassed,
			node.Snapshot.GovernanceFailed,
			node.Snapshot.KPIAverageUtilization*100,
			node.Snapshot.KPIAverageBaseFeeGwei,
		)
	}

	if report.Consensus {
		fmt.Fprintln(out, "Consensus status: OK (all nodes converged on identical state hash)")
		return
	}

	fmt.Fprintf(out, "Consensus status: FAILED (%d divergences)\n", len(report.Divergences))
	for _, divergence := range report.Divergences {
		fmt.Fprintf(out, "  - %s\n", divergence)
	}
}

func shortHash(hash string) string {
	if len(hash) <= 12 {
		return hash
	}
	return hash[:12]
}
