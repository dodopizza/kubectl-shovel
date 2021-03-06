package flags

import (
	"github.com/spf13/pflag"
)

type GCDumpFlagSet struct {
	Timeout Timeout
	dt      *DotnetToolsFlagSet

	flagSet *pflag.FlagSet
}

func NewGCDumpFlagSet() *GCDumpFlagSet {
	return &GCDumpFlagSet{
		Timeout: 30,
		dt:      NewDotnetToolsFlagSet(),
	}
}

func (gc *GCDumpFlagSet) Parse() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("dotnet-gcdump", pflag.ExitOnError)
	flagSet.AddFlagSet(gc.dt.Parse())
	flagSet.Var(
		&gc.Timeout,
		gc.Timeout.Type(),
		gc.Timeout.Description(),
	)

	gc.flagSet = flagSet
	return flagSet
}

func (gc *GCDumpFlagSet) Args() []string {
	args := gc.dt.Args()
	if gc.flagSet.Changed(gc.Timeout.Type()) {
		args = append(
			args,
			FlagToArg(&gc.Timeout)...,
		)
	}
	return args
}
