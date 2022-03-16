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
	descriptionTemplate = "This subcommand will run %s tool for running in k8s application.\n" +
		"Result will be saved locally (or on host) so you'll be able to analyze it with appropriate instruments.\n" +
		"Tool specific additional arguments are also supported.\n" +
		"You can find more info about this tool by the following links:\n\n" +
		"\t* %s\n" +
		"\t* %s"
)

// NewGCDumpCommand return command that start dumper with dotnet-gcdump tool
func NewGCDumpCommand() *cobra.Command {
	builder := NewCommandBuilder(flags.NewDotnetGCDump)
	return builder.Build(
		"Get dotnet-gcdump results",
		fmt.Sprintf(descriptionTemplate,
			builder.Tool(),
			"https://devblogs.microsoft.com/dotnet/collecting-and-analyzing-memory-dumps",
			"https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-gcdump"),
		fmt.Sprintf(examplesTemplate, builder.Tool()),
	)
}

// NewTraceCommand return command that start dumper with dotnet-trace tool
func NewTraceCommand() *cobra.Command {
	builder := NewCommandBuilder(flags.NewDotnetTrace)
	return builder.Build(
		"Get dotnet-trace results",
		fmt.Sprintf(descriptionTemplate,
			builder.Tool(),
			"https://github.com/dotnet/diagnostics/blob/master/documentation/dotnet-trace-instructions.md",
			"https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-trace"),
		fmt.Sprintf(examplesTemplate, builder.Tool()),
	)
}

// NewDumpCommand return command that start dumper with dotnet-dump tool
func NewDumpCommand() *cobra.Command {
	builder := NewCommandBuilder(flags.NewDotnetDump)
	return builder.Build(
		"Get dotnet-dump results",
		fmt.Sprintf(descriptionTemplate,
			builder.Tool(),
			"https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-dump",
			"https://docs.microsoft.com/en-us/dotnet/core/diagnostics/debug-linux-dumps"),
		fmt.Sprintf(examplesTemplate, builder.Tool()),
	)
}

// NewCoreDumpCommand return command that start full process dump with createdump tool
func NewCoreDumpCommand() *cobra.Command {
	builder := NewCommandBuilder(flags.NewCoreDump)
	return builder.Build(
		"Get full process dump results",
		fmt.Sprintf(descriptionTemplate,
			builder.Tool(),
			"https://docs.microsoft.com/en-us/dotnet/core/diagnostics/debug-linux-dumps#core-dumps-with-createdump",
			"https://github.com/dotnet/runtime/blob/main/docs/design/coreclr/botr/xplat-minidump-generation.md"),
		fmt.Sprintf(examplesTemplate, builder.Tool()))
}
