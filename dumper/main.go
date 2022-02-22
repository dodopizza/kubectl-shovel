package main

import (
	"os"

	"github.com/dodopizza/kubectl-shovel/dumper/cmd"
)

func main() {
	rootCmd := cmd.NewDumperCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
