package cmd

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	managedDumpToolName = "dump"
)

type managedDumpOptions struct {
	*commonOptions
	*flags.ManagedDumpFlagSet
}

func newManagedDumpCommand() *cobra.Command {
	options := &managedDumpOptions{
		commonOptions: &commonOptions{},
	}
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", managedDumpToolName),
		Short: "Get dotnet-dump results",
		Long: "This subcommand will run dotnet-dump tool for running in k8s application.\n" +
			"Result will be saved locally so you'll be able to analyze it with appropriate tools.\n" +
			"You can find more info about dotnet-dump tool by the following links:\n\n" +
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-dump\n" +
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/debug-linux-dumps\n",
		Example: fmt.Sprintf(examplesTemplate, managedDumpToolName),
		RunE: func(*cobra.Command, []string) error {
			return options.run()
		},
	}

	cmd.
		Flags().
		AddFlagSet(
			options.parse(),
		)

	return cmd
}

func (options *managedDumpOptions) parse() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(managedDumpToolName, pflag.ExitOnError)
	flagSet.AddFlagSet(options.commonOptions.newCommonFlags(managedDumpToolName))

	options.ManagedDumpFlagSet = flags.NewManagedDumpFlagSet()
	flagSet.AddFlagSet(options.ManagedDumpFlagSet.Parse())

	return flagSet
}

func (options *managedDumpOptions) run() error {
	return run(
		options.commonOptions,
		managedDumpToolName,
		options.Args()...,
	)
}
