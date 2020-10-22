package main

import (
	"os"

	"github.com/dodopizza/kubectl-shovel/cli/cmd"
	"github.com/spf13/pflag"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-shovel", pflag.ExitOnError)
	pflag.CommandLine = flags

	rootCmd := cmd.NewShovelCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
