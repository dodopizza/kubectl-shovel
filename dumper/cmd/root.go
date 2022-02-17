package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
)

var (
	containerInfo = kubernetes.ContainerInfo{}
)

func NewDumperCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "dumper",
		Short:             "Tool to gather diagnostic information from dotnet process",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	cmd.
		PersistentFlags().
		StringVar(
			&containerInfo.ID,
			"container-id",
			containerInfo.ID,
			"Container ID to run tool for",
		)
	cmd.
		PersistentFlags().
		StringVar(
			&containerInfo.Runtime,
			"container-runtime",
			containerInfo.Runtime,
			"Container Runtime to run tool for",
		)
	_ = cmd.MarkPersistentFlagRequired("container-id")
	_ = cmd.MarkPersistentFlagRequired("container-runtime")

	cmd.AddCommand(NewGCDumpCommand())
	cmd.AddCommand(NewTraceCommand())
	cmd.AddCommand(NewDumpCommand())

	return cmd
}
