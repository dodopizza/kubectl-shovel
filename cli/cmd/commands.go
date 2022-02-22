package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dodopizza/kubectl-shovel/internal/flags"
)

var (
	examplesTemplate = "The only required flag is `--pod-name`. So you can use it like this:\n\n" +
		"\tkubectl shovel %[1]s --pod-name my-app-65c4fc589c-gznql\n\n" +
		"Use `-o`/`--output` to define name of dump file:\n\n" +
		"\tkubectl shovel %[1]s --pod-name my-app-65c4fc589c-gznql -o ./myapp.%[1]s\n\n" +
		"Also use `-n`/`--namespace` if your pod is not in current context's namespace:\n\n" +
		"\tkubectl shovel %[1]s --pod-name my-app-65c4fc589c-gznql -n default"
)

// NewGCDumpCommand return command that start dumper with dotnet-gcdump tool
func NewGCDumpCommand() *cobra.Command {
	builder := NewCommandBuilder(flags.NewDotnetGCDump)
	return builder.Build(
		"Get dotnet-gcdump results",
		"This subcommand will run dotnet-gcdump tool for running in k8s application.\n"+
			"Result will be saved locally so you'll be able to analyze it with appropriate tools.\n"+
			"You can find more info about dotnet-gcdump tool by the following links:\n\n"+
			"\t* https://devblogs.microsoft.com/dotnet/collecting-and-analyzing-memory-dumps/\n"+
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-gcdump",
		fmt.Sprintf(examplesTemplate, builder.Tool()),
	)
}

// NewTraceCommand return command that start dumper with dotnet-trace tool
func NewTraceCommand() *cobra.Command {
	builder := NewCommandBuilder(flags.NewDotnetTrace)
	return builder.Build(
		"Get dotnet-trace results",
		"This subcommand will capture runtime events with dotnet-trace tool for running in k8s application.\n"+
			"Result will be saved locally in nettrace format so you'll be able to convert it and analyze with appropriate tools.\n"+
			"You can find more info about dotnet-trace tool by the following links:\n\n"+
			"\t* https://github.com/dotnet/diagnostics/blob/master/documentation/dotnet-trace-instructions.md\n"+
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-trace",
		fmt.Sprintf(examplesTemplate, builder.Tool())+"\n\n"+
			"Use `--duration` to define duration of trace to 30 seconds:\n\n"+
			"\tkubectl shovel trace --pod-name my-app-65c4fc589c-gznql -o ./myapp.trace --duration 30s\n\n"+
			"Use `--format` to specify Speedscope format:\n\n"+
			"\tkubectl shovel trace --pod-name my-app-65c4fc589c-gznql -o ./myapp.trace --format Speedscope\n\n"+
			"And then you can analyze it with https://www.speedscope.app/\n"+
			"Or convert any other format to speedscope format with:\n\n"+
			"\tdotnet trace convert myapp.trace --format Speedscope",
	)
}

// NewDumpCommand return command that start dumper with dotnet-dump tool
func NewDumpCommand() *cobra.Command {
	builder := NewCommandBuilder(flags.NewDotnetDump)
	return builder.Build(
		"Get dotnet-dump results",
		"This subcommand will run dotnet-dump tool for running in k8s application.\n"+
			"Result will be saved locally so you'll be able to analyze it with appropriate tools.\n"+
			"You can find more info about dotnet-dump tool by the following links:\n\n"+
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-dump\n"+
			"\t* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/debug-linux-dumps\n",
		fmt.Sprintf(examplesTemplate, builder.Tool()),
	)
}
