package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dodopizza/kubectl-shovel/internal/globals"
)

// NewShovelCommand initialize Shovel command
func NewShovelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               globals.PluginName,
		Short:             "Get diagnostics from running in k8s dotnet application",
		Version:           globals.GetVersion(),
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	cmd.AddCommand(NewDocCommand())
	cmd.AddCommand(NewVersionCommand())
	cmd.AddCommand(NewGCDumpCommand())
	cmd.AddCommand(NewTraceCommand())
	cmd.AddCommand(NewDumpCommand())
	cmd.AddCommand(NewCreateDumpCommand())

	return cmd
}
