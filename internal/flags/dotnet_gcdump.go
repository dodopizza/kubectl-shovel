package flags

import (
	"github.com/spf13/pflag"

	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
)

type DotnetGCDump struct {
	*DotnetToolProperties

	Timeout types.Timeout

	flagSet *pflag.FlagSet
}

func NewDotnetGCDump() DotnetTool {
	return &DotnetGCDump{
		DotnetToolProperties: NewDotnetToolProperties(),
		Timeout:              30,
	}
}

func (gc *DotnetGCDump) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(gc.BinaryName(), pflag.ExitOnError)
	flagSet.AddFlagSet(gc.DotnetToolProperties.GetFlags())
	flagSet.Var(
		&gc.Timeout,
		gc.Timeout.Type(),
		gc.Timeout.Description(),
	)

	gc.flagSet = flagSet
	return flagSet
}

func (gc *DotnetGCDump) FormatArgs(args *Args) {
	args.AppendFrom(gc.DotnetToolProperties)
	if gc.flagSet.Changed(gc.Timeout.Type()) {
		args.AppendFlag(&gc.Timeout)
	}
}

func (*DotnetGCDump) BinaryName() string {
	return "dotnet-gcdump"
}

func (*DotnetGCDump) ToolName() string {
	return "gcdump"
}
