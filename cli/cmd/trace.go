package cmd

import (
	"fmt"

	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
)

func newTraceCommand() *cobra.Command {
	options := NewDiagnosticToolOptions("trace", flags.NewDotnetTrace)
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", options.Tool),
		Short: "Get dotnet-trace results",
		Long: "This subcommand will capture runtime events with dotnet-trace tool for running in k8s appplication.\n" +
			"Result will be saved locally in nettrace format so you'll be able to convert it and analyze with appropriate tools.\n" +
			"You can find more info about dotnet-trace tool by the following links:\n\n" +
			"\t* https://github.com/dotnet/diagnostics/blob/master/documentation/dotnet-trace-instructions.md\n" +
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-trace",
		Example: fmt.Sprintf(examplesTemplate, options.Tool) + "\n\n" +
			"Use `--duration` to define duration of trace to 30 seconds:\n\n" +
			"\tkubectl shovel trace --pod-name my-app-65c4fc589c-gznql -o ./myapp.trace --duration 30s\n\n" +
			"Use `--format` to specify Speedscope format:\n\n" +
			"\tkubectl shovel trace --pod-name my-app-65c4fc589c-gznql -o ./myapp.trace --format Speedscope\n\n" +
			"And then you can analyze it with https://www.speedscope.app/\n" +
			"Or convert any other format to speedscope format with:\n\n" +
			"\tdotnet trace convert myapp.trace --format Speedscope",
		RunE: func(*cobra.Command, []string) error {
			return options.Run()
		},
	}

	cmd.Flags().AddFlagSet(options.Parse())
	return cmd
}
