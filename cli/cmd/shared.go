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

type DiagnosticToolOptions struct {
	CommonOptions *commonOptions
	FlagSet       flags.DotnetTool
	Tool          string
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

func NewDiagnosticToolOptions(tool string, factory flags.DotnetToolFactory) *DiagnosticToolOptions {
	return &DiagnosticToolOptions{
		CommonOptions: &commonOptions{},
		FlagSet:       factory(),
		Tool:          tool,
	}
}

func (dt *DiagnosticToolOptions) Parse() *pflag.FlagSet {
	fs := pflag.NewFlagSet(dt.Tool, pflag.ExitOnError)
	fs.AddFlagSet(dt.CommonOptions.GetFlags(dt.Tool))
	fs.AddFlagSet(dt.FlagSet.GetFlags())
	return fs
}

func (dt *DiagnosticToolOptions) Run() error {
	return launch(dt.CommonOptions, dt.Tool, dt.FlagSet.GetArgs()...)
}
