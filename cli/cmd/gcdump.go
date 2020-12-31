package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	gcDumpToolName = "gcdump"
)

type gcDumpOptions struct {
	*commonOptions
}

func newGCDumpCommand() *cobra.Command {
	options := &gcDumpOptions{
		commonOptions: &commonOptions{},
	}
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", gcDumpToolName),
		Short: "Get dotnet-gcdump results",
		Long: "This subcommand will run dotnet-gcdump tool for running in k8s appplication.\n" +
			"Result will be saved locally so you'll be able to analyze it with appropriate tools.\n" +
			"You can find more info about dotnet-gcdump tool by the following links:\n\n" +
			"\t* https://devblogs.microsoft.com/dotnet/collecting-and-analyzing-memory-dumps/\n" +
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-gcdump",
		Example: fmt.Sprintf(examplesTemplate, gcDumpToolName),
		RunE: func(*cobra.Command, []string) error {
			return options.makeGCDump()
		},
	}

	cmd.
		Flags().
		AddFlagSet(
			options.parseFlags(),
		)

	return cmd
}

func (options *gcDumpOptions) parseFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(gcDumpToolName, pflag.ExitOnError)
	flags.AddFlagSet(options.commonOptions.newCommonFlags(gcDumpToolName))

	return flags
}

func (options *gcDumpOptions) makeGCDump() error {
	return run(
		options.commonOptions,
		gcDumpToolName,
	)
}
