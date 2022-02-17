package cmd

import (
	"fmt"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/dodopizza/kubectl-shovel/internal/globals"
)

type commonOptions struct {
	image         string
	podName       string
	output        string
	containerName string

	kubeFlags *genericclioptions.ConfigFlags
}

type CommandBuilder struct {
	CommonOptions *commonOptions
	tool          flags.DotnetTool
}

func (options *commonOptions) GetFlags(tool string) *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("common", pflag.ExitOnError)
	flagSet.StringVar(
		&options.podName,
		"pod-name",
		options.podName,
		"Target pod",
	)
	_ = cobra.MarkFlagRequired(flagSet, "pod-name")

	flagSet.StringVarP(
		&options.containerName,
		"container",
		"c",
		options.containerName,
		"Target container in pod. Required if pod run multiple containers",
	)

	flagSet.StringVarP(
		&options.output,
		"output",
		"o",
		fmt.Sprintf(
			"./output.%s",
			tool,
		),
		"Output file",
	)

	flagSet.StringVar(
		&options.image,
		"image",
		globals.GetDumperImage(),
		"Image of dumper to use for job",
	)
	options.kubeFlags = genericclioptions.NewConfigFlags(false)
	options.kubeFlags.AddFlags(flagSet)

	return flagSet
}

func NewCommandBuilder(factory flags.DotnetToolFactory) *CommandBuilder {
	return &CommandBuilder{
		CommonOptions: &commonOptions{},
		tool:          factory(),
	}
}

func (dt *CommandBuilder) Build(short, long, example string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("%s [flags]", dt.tool.ToolName()),
		Short:   short,
		Long:    long,
		Example: example,
		RunE: func(*cobra.Command, []string) error {
			return dt.run()
		},
	}
	cmd.Flags().AddFlagSet(dt.parse())

	return cmd
}

func (dt *CommandBuilder) Tool() string {
	return dt.tool.ToolName()
}

func (dt *CommandBuilder) parse() *pflag.FlagSet {
	fs := pflag.NewFlagSet(dt.tool.ToolName(), pflag.ExitOnError)
	fs.AddFlagSet(dt.CommonOptions.GetFlags(dt.tool.ToolName()))
	fs.AddFlagSet(dt.tool.GetFlags())
	return fs
}

func (dt *CommandBuilder) run() error {
	return launch(dt.CommonOptions, dt.tool.ToolName(), dt.tool.GetArgs()...)
}
