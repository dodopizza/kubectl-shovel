package cmd

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CommonOptions represents container info
type CommonOptions struct {
	ContainerID       string
	ContainerRuntime  string
	ContainerName     string
	PodName           string
	PodNamespace      string
	StoreOutputOnHost bool
}

// CommandBuilder represents logic for building and running tools
type CommandBuilder struct {
	CommonOptions *CommonOptions
	tool          flags.DotnetTool
}

// GetFlags return FlagSet that describes options for container selection
func (options *CommonOptions) GetFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("common", pflag.ExitOnError)
	fs.StringVar(
		&options.ContainerID,
		"container-id",
		options.ContainerID,
		"Container ID to run tool for",
	)
	_ = cobra.MarkFlagRequired(fs, "container-id")

	fs.StringVar(
		&options.ContainerRuntime,
		"container-runtime",
		options.ContainerRuntime,
		"Container Runtime to run tool for",
	)
	_ = cobra.MarkFlagRequired(fs, "container-runtime")

	fs.StringVar(
		&options.ContainerName,
		"container-name",
		options.ContainerName,
		"Container name to run tool for",
	)

	fs.StringVar(
		&options.ContainerRuntime,
		"pod-name",
		options.PodName,
		"Pod name to run tool for",
	)

	fs.StringVar(
		&options.ContainerRuntime,
		"pod-namespace",
		options.PodName,
		"Pod namespace to run tool for",
	)

	fs.BoolVar(
		&options.StoreOutputOnHost,
		"store-output-on-host",
		options.StoreOutputOnHost,
		"Flag, indicating that output should be stored on host /tmp folder",
	)

	return fs
}

// NewCommandBuilder returns options with specified tool name
func NewCommandBuilder(options *CommonOptions, factory flags.DotnetToolFactory) *CommandBuilder {
	return &CommandBuilder{
		tool:          factory(),
		CommonOptions: options,
	}
}

// Build will build command *cobra.Command from options
func (cb *CommandBuilder) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", cb.tool.ToolName()),
		Args:  cobra.NoArgs,
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cb.launch()
		},
	}

	cmd.PersistentFlags().AddFlagSet(cb.flags())
	return cmd
}

func (cb *CommandBuilder) flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet(cb.tool.ToolName(), pflag.ExitOnError)
	fs.AddFlagSet(cb.tool.GetFlags())
	return fs
}
