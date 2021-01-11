package cmd

import (
	"os"

	"github.com/dodopizza/kubectl-shovel/internal/events"
	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
	"github.com/spf13/cobra"
)

var (
	containerInfo = kubernetes.ContainerInfo{}
)

var rootCmd = &cobra.Command{
	Use:               "dumper",
	Short:             "Tool to gather diagnostic information from dotnet process",
	SilenceUsage:      true,
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
			&containerInfo.ID,
			"container-id",
			containerInfo.ID,
			"Container ID to run tool for",
		)
	rootCmd.
		PersistentFlags().
		StringVar(
			&containerInfo.Runtime,
			"container-runtime",
			containerInfo.Runtime,
			"Container Runtime to run tool for",
		)
	err := rootCmd.MarkPersistentFlagRequired("container-runtime")
	if err != nil {
		return err
	}

	rootCmd.AddCommand(newGCDumpCommand())
	rootCmd.AddCommand(newTraceCommand())

	return nil
}
