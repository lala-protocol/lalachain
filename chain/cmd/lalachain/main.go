package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lala-protocol/lalachain/chain/app"
	"github.com/lala-protocol/lalachain/chain/cosmosruntime"
)

func main() {
	epochs := flag.Int("epochs", 30, "Number of epochs to run")
	seed := flag.Int64("seed", 42, "Random seed for reproducible simulation")
	runtime := flag.String("runtime", "prototype", "Runtime mode: prototype|cosmos")
	network := flag.String("network", "single", "Network mode: single|multi")
	phase := flag.String("phase", "1", "Execution phase: 1|2 (or phase1|phase2)")
	validators := flag.Int("validators", 10, "Validator count for phase 2 mode (4-10)")
	nodes := flag.Int("nodes", 0, "Node count for multi network mode (default: validators)")
	flag.Parse()

	selectedRuntime := strings.ToLower(strings.TrimSpace(*runtime))
	selectedNetwork := strings.ToLower(strings.TrimSpace(*network))
	selected := strings.ToLower(strings.TrimSpace(*phase))
	switch selectedRuntime {
	case "prototype":
		if selectedNetwork != "single" {
			fmt.Fprintf(os.Stderr, "prototype runtime only supports --network single\n")
			os.Exit(2)
		}
		prototype := app.NewPrototype(*seed)
		switch selected {
		case "1", "phase1":
			prototype = app.NewPrototype(*seed)
		case "2", "phase2":
			prototype = app.NewPhase2Prototype(*seed, *validators)
		default:
			fmt.Fprintf(os.Stderr, "invalid phase %q (supported: 1, 2, phase1, phase2)\n", *phase)
			os.Exit(2)
		}

		if err := prototype.Run(*epochs, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "prototype run failed: %v\n", err)
			os.Exit(1)
		}
	case "cosmos":
		switch selectedNetwork {
		case "single":
			runner, err := cosmosruntime.NewRunner(*seed, selected, *validators)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(2)
			}
			if err := runner.Run(sdk.Context{}, *epochs, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "cosmos scaffold run failed: %v\n", err)
				os.Exit(1)
			}
		case "multi":
			report, err := cosmosruntime.RunLocalTestnet(cosmosruntime.LocalTestnetConfig{
				Phase:      selected,
				Epochs:     *epochs,
				Seed:       *seed,
				Validators: *validators,
				Nodes:      *nodes,
			}, os.Stdout)
			if err != nil {
				fmt.Fprintf(os.Stderr, "local cosmos testnet failed: %v\n", err)
				os.Exit(1)
			}
			if !report.Consensus {
				fmt.Fprintln(os.Stderr, "local cosmos testnet failed consensus checks")
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "invalid network %q (supported: single, multi)\n", *network)
			os.Exit(2)
		}
	default:
		fmt.Fprintf(os.Stderr, "invalid runtime %q (supported: prototype, cosmos)\n", *runtime)
		os.Exit(2)
	}
}
