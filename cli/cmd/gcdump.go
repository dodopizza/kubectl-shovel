package cmd

import (
	"fmt"

	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
)

func newGCDumpCommand() *cobra.Command {
	options := NewDiagnosticToolOptions("gcdump", flags.NewGCDumpFlagSet)
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", options.Tool),
		Short: "Get dotnet-gcdump results",
		Long: "This subcommand will run dotnet-gcdump tool for running in k8s appplication.\n" +
			"Result will be saved locally so you'll be able to analyze it with appropriate tools.\n" +
			"You can find more info about dotnet-gcdump tool by the following links:\n\n" +
			"\t* https://devblogs.microsoft.com/dotnet/collecting-and-analyzing-memory-dumps/\n" +
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-gcdump",
		Example: fmt.Sprintf(examplesTemplate, options.Tool),
		RunE: func(*cobra.Command, []string) error {
			return options.Run()
		},
	}

	cmd.Flags().AddFlagSet(options.Parse())
	return cmd
}
