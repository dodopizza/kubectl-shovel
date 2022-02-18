package cmd

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/dodopizza/kubectl-shovel/internal/globals"
)

// CommonOptions contains generic arguments for cli
type CommonOptions struct {
	Container string
	Image     string
	Pod       string
	Output    string

	kube *genericclioptions.ConfigFlags
}

// CommandBuilder represents logic for building and running tools
type CommandBuilder struct {
	CommonOptions *CommonOptions
	tool          flags.DotnetTool
}

// GetFlags return FlagSet that describes generic options
func (options *CommonOptions) GetFlags(tool string) *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("common", pflag.ExitOnError)
	flagSet.StringVar(
		&options.Pod,
		"pod-name",
		options.Pod,
		"Target pod",
	)
	_ = cobra.MarkFlagRequired(flagSet, "pod-name")

	flagSet.StringVarP(
		&options.Container,
		"container",
		"c",
		options.Container,
		"Target container in pod. Required if pod run multiple containers",
	)

	flagSet.StringVarP(
		&options.Output,
		"output",
		"o",
		fmt.Sprintf(
			"./output.%s",
			tool,
		),
		"Output file",
	)

	flagSet.StringVar(
		&options.Image,
		"image",
		globals.GetDumperImage(),
		"Image of dumper to use for job",
	)
	options.kube = genericclioptions.NewConfigFlags(false)
	options.kube.AddFlags(flagSet)

	return flagSet
}

func NewCommandBuilder(factory flags.DotnetToolFactory) *CommandBuilder {
	return &CommandBuilder{
		CommonOptions: &CommonOptions{},
		tool:          factory(),
	}
}

// Build will build and returns *cobra.Command for specific tool
func (cb *CommandBuilder) Build(short, long, example string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("%s [flags]", cb.tool.ToolName()),
		Short:   short,
		Long:    long,
		Example: example,
		RunE: func(*cobra.Command, []string) error {
			return cb.launch()
		},
	}
	cmd.Flags().AddFlagSet(cb.parse())

	return cmd
}

// Tool returns tool name that contains CommandBuilder
func (cb *CommandBuilder) Tool() string {
	return cb.tool.ToolName()
}

func (cb *CommandBuilder) parse() *pflag.FlagSet {
	fs := pflag.NewFlagSet(cb.tool.ToolName(), pflag.ExitOnError)
	fs.AddFlagSet(cb.CommonOptions.GetFlags(cb.tool.ToolName()))
	fs.AddFlagSet(cb.tool.GetFlags())
	return fs
}
