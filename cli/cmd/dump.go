package cmd

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	dumpToolName = "dump"
)

type dumpOptions struct {
	*commonOptions
	*flags.DumpFlagSet
}

func newDumpCommand() *cobra.Command {
	options := &dumpOptions{
		commonOptions: &commonOptions{},
	}
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", dumpToolName),
		Short: "Get dotnet-dump results",
		Long: "This subcommand will run dotnet-dump tool for running in k8s application.\n" +
			"Result will be saved locally so you'll be able to analyze it with appropriate tools.\n" +
			"You can find more info about dotnet-dump tool by the following links:\n\n" +
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-dump\n" +
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/debug-linux-dumps\n",
		Example: fmt.Sprintf(examplesTemplate, dumpToolName),
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

func (options *dumpOptions) parse() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(dumpToolName, pflag.ExitOnError)
	flagSet.AddFlagSet(options.commonOptions.newCommonFlags(dumpToolName))

	options.DumpFlagSet = flags.NewDumpFlagSet()
	flagSet.AddFlagSet(options.DumpFlagSet.Parse())

	return flagSet
}

func (options *dumpOptions) run() error {
	return run(
		options.commonOptions,
		dumpToolName,
		options.Args()...,
	)
}
