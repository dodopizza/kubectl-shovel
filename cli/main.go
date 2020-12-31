package main

import (
	"os"

	"github.com/dodopizza/kubectl-shovel/cli/cmd"
)

func main() {
	rootCmd := cmd.NewShovelCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
