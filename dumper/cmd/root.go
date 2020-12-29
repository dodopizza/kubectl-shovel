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
	if err := rootCmd.Execute(); err != nil {
		events.NewEvent(events.Error, err.Error())
		os.Exit(1)
	}
}

func init() {
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
	cobra.MarkFlagFilename(rootCmd.PersistentFlags(), "container-runtime")

	rootCmd.AddCommand(newGCDumpCommand())
	rootCmd.AddCommand(newTraceCommand())
}
