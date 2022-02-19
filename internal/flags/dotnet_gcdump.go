package flags

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
	"github.com/spf13/pflag"
)

type DotnetGCDump struct {
	Timeout types.Timeout
	dt      *DotnetToolProperties

	flagSet *pflag.FlagSet
}

func NewDotnetGCDump() DotnetTool {
	return &DotnetGCDump{
		Timeout: 30,
		dt:      NewDotnetToolProperties(),
	}
}

func (gc *DotnetGCDump) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(gc.BinaryName(), pflag.ExitOnError)
	flagSet.AddFlagSet(gc.dt.GetFlags())
	flagSet.Var(
		&gc.Timeout,
		gc.Timeout.Type(),
		gc.Timeout.Description(),
	)

	gc.flagSet = flagSet
	return flagSet
}

func (gc *DotnetGCDump) GetArgs() []string {
	args := gc.dt.GetArgs()
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

func (gc *DotnetGCDump) GetProperties() *DotnetToolProperties {
	return gc.dt
}
