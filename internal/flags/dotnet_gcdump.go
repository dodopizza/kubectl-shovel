package flags

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
	"github.com/spf13/pflag"
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

func (gc *DotnetGCDump) FormatArgs() []string {
	args := gc.DotnetToolProperties.FormatArgs()
	if gc.flagSet.Changed(gc.Timeout.Type()) {
		args = append(
			args,
			FlagToArg(&gc.Timeout)...,
		)
	}
	return args
}

func (gc *DotnetGCDump) BinaryName() string {
	return "dotnet-gcdump"
}

func (gc *DotnetGCDump) ToolName() string {
	return "gcdump"
}

func (gc *DotnetGCDump) GetProperties() DotnetToolFlags {
	return gc.DotnetToolProperties
}
