package flags

import (
	"strconv"

	"github.com/spf13/pflag"
)

type DotnetToolSharedOptions struct {
	Action    string
	Output    string
	ProcessID int
}

func NewDotnetToolSharedOptions() *DotnetToolSharedOptions {
	return &DotnetToolSharedOptions{
		ProcessID: 1,
	}
}

func (dt *DotnetToolSharedOptions) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("dotnet-tools", pflag.ExitOnError)
	flagSet.IntVarP(
		&dt.ProcessID,
		"process-id",
		"p",
		dt.ProcessID,
		"The process ID to collect the trace from",
	)
	flagSet.StringVarP(
		&dt.Output,
		"output",
		"o",
		dt.Output,
		"Output file",
	)
	return flagSet
}

func (dt *DotnetToolSharedOptions) FormatArgs(args *Args, t FormatArgsType) {
	// append command (collect) for binary execution
	if t == FormatArgsTypeBinary {
		args.AppendRaw("collect")
	}

	args.Append("process-id", strconv.Itoa(dt.ProcessID))

	if dt.Output != "" {
		args.Append("output", dt.Output)
	}
}

func (dt *DotnetToolSharedOptions) SetOutput(output string) {
	dt.Output = output
}

func (dt *DotnetToolSharedOptions) SetProcessID(id int) {
	dt.ProcessID = id
}
