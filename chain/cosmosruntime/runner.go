package cosmosruntime

import (
	"errors"
	"fmt"
	"io"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lala-protocol/lalachain/chain/app"
	"github.com/lala-protocol/lalachain/chain/cosmosapp"
)

// Runner executes the adaptation loop through a Cosmos SDK-typed boundary.
// This keeps module wiring compatible with future full SDK app integration.
type Runner struct {
	app   *cosmosapp.App
	phase app.Phase
}

func NewRunner(seed int64, phase string, validatorCount int) (*Runner, error) {
	selected := strings.ToLower(strings.TrimSpace(phase))
	switch selected {
	case "1", "phase1":
		cosmosApp, err := cosmosapp.New(seed, app.Phase1, validatorCount)
		if err != nil {
			return nil, err
		}
		return &Runner{
			app:   cosmosApp,
			phase: app.Phase1,
		}, nil
	case "2", "phase2":
		cosmosApp, err := cosmosapp.New(seed, app.Phase2, validatorCount)
		if err != nil {
			return nil, err
		}
		return &Runner{
			app:   cosmosApp,
			phase: app.Phase2,
		}, nil
	default:
		return nil, fmt.Errorf("invalid phase %q (supported: 1, 2, phase1, phase2)", phase)
	}
}

func (r *Runner) Run(ctx sdk.Context, epochs int, out io.Writer) error {
	if r == nil || r.app == nil {
		return errors.New("cosmos runtime is not initialized")
	}
	if epochs <= 0 {
		return errors.New("epochs must be > 0")
	}

	// Context is intentionally accepted even though state is still in-memory,
	// so callers can switch to SDK context-backed state without API churn.
	_ = ctx

	fmt.Fprintln(out, "======================================================================")
	fmt.Fprintln(out, " Cosmos SDK Scaffold Runtime")
	fmt.Fprintln(out, "======================================================================")
	fmt.Fprintf(out, "Phase: %s\n\n", r.phase)

	return r.app.Run(epochs, out)
}

func (r *Runner) Snapshot() (cosmosapp.Snapshot, bool) {
	if r == nil || r.app == nil {
		return cosmosapp.Snapshot{}, false
	}
	return r.app.Snapshot()
}

func (r *Runner) StateHash() string {
	if r == nil || r.app == nil {
		return ""
	}
	return r.app.StateHash()
}
