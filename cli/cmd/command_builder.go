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

func (cb *CommandBuilder) Build(short, long, example string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("%s [flags]", cb.tool.ToolName()),
		Short:   short,
		Long:    long,
		Example: example,
		RunE: func(*cobra.Command, []string) error {
			return cb.run()
		},
	}
	cmd.Flags().AddFlagSet(cb.parse())

	return cmd
}

func (cb *CommandBuilder) Tool() string {
	return cb.tool.ToolName()
}

func (cb *CommandBuilder) parse() *pflag.FlagSet {
	fs := pflag.NewFlagSet(cb.tool.ToolName(), pflag.ExitOnError)
	fs.AddFlagSet(cb.CommonOptions.GetFlags(cb.tool.ToolName()))
	fs.AddFlagSet(cb.tool.GetFlags())
	return fs
}

func (cb *CommandBuilder) run() error {
	return launch(cb.CommonOptions, cb.tool.ToolName(), cb.tool.GetArgs()...)
}
