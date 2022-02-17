package main

import (
	"github.com/dodopizza/kubectl-shovel/dumper/cmd"
	"github.com/dodopizza/kubectl-shovel/internal/events"
	"os"
)

func main() {
	rootCmd := cmd.NewDumperCommand()
	if err := rootCmd.Execute(); err != nil {
		events.NewEvent(events.Error, err.Error())
		os.Exit(1)
	}
}
