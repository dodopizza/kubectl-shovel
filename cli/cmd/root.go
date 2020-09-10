package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pluginName      = "kubectl-shovel"
	dumperImageName = "dodoreg.azurecr.io/dumper"
)

// NewShovelCommand initialize Shovel command
func NewShovelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pluginName,
		Short: "Get diagnostics from running in k8s dotnet application",
	}

	cmd.AddCommand(newDocCmd())
	cmd.AddCommand(newVersionCmd())

	cmd.AddCommand(newGCDumpCommand())
	cmd.AddCommand(newTraceCommand())

	return cmd
}
