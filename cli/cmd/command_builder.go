package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/dodopizza/kubectl-shovel/internal/globals"
	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
)

// CommonOptions contains shared arguments for cli commands
type CommonOptions struct {
	Container         string
	Image             string
	Pod               string
	Output            string
	OutputHostPath    string
	LimitCPU          string
	LimitMemory       string
	StoreOutputOnHost bool

	kubeConfig *genericclioptions.ConfigFlags
}

// CommandBuilder represents logic for building and running tools
type CommandBuilder struct {
	CommonOptions *CommonOptions
	tool          flags.DotnetTool
	kube          *kubernetes.Client
}

// GetFlags return FlagSet that describes generic options
func (options *CommonOptions) GetFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("common", pflag.ExitOnError)
	fs.StringVarP(
		&options.Container,
		"container",
		"c",
		options.Container,
		"Target container in pod. Required if pod run multiple containers. Will look in init containers if not found in regular containers",
	)
	fs.StringVar(
		&options.Image,
		"image",
		options.Image,
		"Image of dumper to use for job",
	)
	fs.StringVar(
		&options.Pod,
		"pod-name",
		options.Pod,
		"Target pod",
	)
	_ = cobra.MarkFlagRequired(fs, "pod-name")
	fs.StringVarP(
		&options.Output,
		"output",
		"o",
		options.Output,
		"Output file",
	)
	fs.StringVar(
		&options.OutputHostPath,
		"output-host-path",
		options.OutputHostPath,
		"Host folder, where will be stored artifact",
	)
	fs.StringVar(
		&options.LimitCPU,
		"limit-cpu",
		options.LimitCPU,
		"Limit maximal consumptions cpu for the executing job",
	)
	fs.StringVar(
		&options.LimitMemory,
		"limit-memory",
		options.LimitMemory,
		"Limit maximal consumptions memory for the executing job",
	)
	fs.BoolVarP(
		&options.StoreOutputOnHost,
		"store-output-on-host",
		"t",
		options.StoreOutputOnHost,
		"Store output on node instead of downloading it locally")

	options.kubeConfig = genericclioptions.NewConfigFlags(false)
	options.kubeConfig.AddFlags(fs)

	return fs
}

// NewCommandBuilder returns *CommandBuilder instance with specified factory flags.DotnetToolFactory
// that responsible for creation of any available flags.DotnetTool
func NewCommandBuilder(factory flags.DotnetToolFactory) *CommandBuilder {
	tool := factory()

	return &CommandBuilder{
		CommonOptions: &CommonOptions{
			Image:          globals.GetDumperImage(),
			Output:         fmt.Sprintf("./output.%s", tool.ToolName()),
			OutputHostPath: fmt.Sprintf("%s/%s", globals.PathTmpFolder, globals.PluginName),
		},
		tool: tool,
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
	fs.AddFlagSet(cb.CommonOptions.GetFlags())
	fs.AddFlagSet(cb.tool.GetFlags())
	return fs
}
