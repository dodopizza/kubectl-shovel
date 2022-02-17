package cmd

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
)

func newDumpCommand() *cobra.Command {
	options := NewDiagnosticToolOptions("dump", flags.NewDumpFlagSet)
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", options.Tool),
		Short: "Get dotnet-dump results",
		Long: "This subcommand will run dotnet-dump tool for running in k8s application.\n" +
			"Result will be saved locally so you'll be able to analyze it with appropriate tools.\n" +
			"You can find more info about dotnet-dump tool by the following links:\n\n" +
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-dump\n" +
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/debug-linux-dumps\n",
		Example: fmt.Sprintf(examplesTemplate, options.Tool),
		RunE: func(*cobra.Command, []string) error {
			return options.Run()
		},
	}

	cmd.Flags().AddFlagSet(options.Parse())
	return cmd
}
