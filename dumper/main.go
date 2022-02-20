package main

import (
	"github.com/dodopizza/kubectl-shovel/dumper/cmd"
	"os"
)

func main() {
	rootCmd := cmd.NewDumperCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
