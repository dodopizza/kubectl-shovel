package flags

import (
	"strconv"

	"github.com/spf13/pflag"
)

type DotnetToolsFlagSet struct {
	ProcessID int
}

func NewDotnetToolsFlagSet() *DotnetToolsFlagSet {
	return &DotnetToolsFlagSet{
		ProcessID: 1,
	}
}

func (dt *DotnetToolsFlagSet) Parse() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("dotnet-tools", pflag.ExitOnError)
	flagSet.IntVarP(
		&dt.ProcessID,
		"process-id",
		"p",
		dt.ProcessID,
		"The process ID to collect the trace from",
	)

	return flagSet
}

func (dt *DotnetToolsFlagSet) Args() []string {
	return []string{
		"--process-id", strconv.Itoa(dt.ProcessID),
	}
}
