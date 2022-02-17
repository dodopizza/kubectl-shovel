package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dodopizza/kubectl-shovel/internal/version"
)

const (
	dumperImageName = "dodopizza/kubectl-shovel-dumper"
)

// NewShovelCommand initialize Shovel command
func NewShovelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "kubectl-shovel",
		Short:             "Get diagnostics from running in k8s dotnet application",
		Version:           version.GetVersion(),
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	cmd.AddCommand(newDocCmd())
	cmd.AddCommand(newVersionCmd())

	cmd.AddCommand(newGCDumpCommand())
	cmd.AddCommand(newTraceCommand())
	cmd.AddCommand(newDumpCommand())

	return cmd
}
