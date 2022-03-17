package flags

import (
	"github.com/spf13/pflag"

	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
)

type gcdump struct {
	*DotnetToolSharedOptions

	Timeout types.Timeout

	flagSet *pflag.FlagSet
}

func NewDotnetGCDump() DotnetTool {
	return &gcdump{
		DotnetToolSharedOptions: NewDotnetToolSharedOptions(),
		Timeout:                 30,
	}
}

func (gc *gcdump) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(gc.BinaryName(), pflag.ExitOnError)
	flagSet.AddFlagSet(gc.DotnetToolSharedOptions.GetFlags())
	flagSet.Var(
		&gc.Timeout,
		gc.Timeout.Type(),
		gc.Timeout.Description(),
	)

	gc.flagSet = flagSet
	return flagSet
}

func (gc *gcdump) FormatArgs(args *Args, t FormatArgsType) {
	gc.DotnetToolSharedOptions.FormatArgs(args, t)
	if gc.flagSet.Changed(gc.Timeout.Type()) {
		args.AppendFlag(&gc.Timeout)
	}
}

func (*gcdump) BinaryName() string {
	return "dotnet-gcdump"
}

func (*gcdump) ToolName() string {
	return "gcdump"
}
