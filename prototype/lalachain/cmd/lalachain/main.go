package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lala-protocol/whitepaper/prototype/lalachain/app"
)

func main() {
	epochs := flag.Int("epochs", 30, "Number of epochs to run")
	seed := flag.Int64("seed", 42, "Random seed for reproducible simulation")
	flag.Parse()

	prototype := app.NewPrototype(*seed)
	if err := prototype.Run(*epochs, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "phase 1 prototype failed: %v\n", err)
		os.Exit(1)
	}
}
