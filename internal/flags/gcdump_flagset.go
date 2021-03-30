package flags

import (
	"strconv"

	"github.com/spf13/pflag"
)

type GCDumpFlagSet struct {
	Timeout int
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
	flagSet.IntVar(
		&gc.Timeout,
		"timeout",
		gc.Timeout,
		"Give up on collecting the GC dump if it takes longer than this many seconds",
	)

	gc.flagSet = flagSet
	return flagSet
}

func (gc *GCDumpFlagSet) Args() []string {
	args := gc.dt.Args()
	if gc.flagSet.Changed("timeout") {
		args = append(
			args,
			[]string{
				"--timeout", strconv.Itoa(gc.Timeout),
			}...,
		)
	}
	return args
}
