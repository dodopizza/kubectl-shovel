package cmd

import (
	"os"

	"github.com/dodopizza/kubectl-shovel/pkg/events"
	"github.com/spf13/cobra"
)

var (
	containerID      = ""
	containerRuntime = ""
)

var rootCmd = &cobra.Command{
	Use:               "dumper",
	Short:             "Tool to gather diagnostic information from dotnet process",
	DisableAutoGenTag: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := initializeRootCmd()
	if err != nil {
		events.NewEvent(events.Error, err.Error())
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		events.NewEvent(events.Error, err.Error())
		os.Exit(1)
	}
}

func initializeRootCmd() error {
	rootCmd.
		PersistentFlags().
		StringVar(
			&containerID,
			"container-id",
			containerID,
			"Container ID to run tool for",
		)
	rootCmd.
		PersistentFlags().
		StringVar(
			&containerRuntime,
			"container-runtime",
			containerRuntime,
			"Container ID to run tool for",
		)
	err := rootCmd.MarkPersistentFlagRequired("container-runtime")
	if err != nil {
		return err
	}

	rootCmd.AddCommand(newGCDumpCommand())
	rootCmd.AddCommand(newTraceCommand())

	return nil
}