package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/lala-protocol/lalachain/chain/app/chain"
	"github.com/lala-protocol/lalachain/chain/cmd/lalachaind/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", chain.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
