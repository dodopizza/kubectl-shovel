package flags

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
	"github.com/spf13/pflag"
)

type GCDumpFlagSet struct {
	Timeout types.Timeout
	dt      *DotnetToolsFlagSet

	flagSet *pflag.FlagSet
}

func NewGCDumpFlagSet() FlagSetContainer {
	return &GCDumpFlagSet{
		Timeout: 30,
		dt:      NewDotnetToolsFlagSet(),
	}
}

func (gc *GCDumpFlagSet) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("dotnet-gcdump", pflag.ExitOnError)
	flagSet.AddFlagSet(gc.dt.GetFlags())
	flagSet.Var(
		&gc.Timeout,
		gc.Timeout.Type(),
		gc.Timeout.Description(),
	)

	gc.flagSet = flagSet
	return flagSet
}

func (gc *GCDumpFlagSet) GetArgs() []string {
	args := gc.dt.GetArgs()
	if gc.flagSet.Changed(gc.Timeout.Type()) {
		args = append(
			args,
			FlagToArg(&gc.Timeout)...,
		)
	}
	return args
}
