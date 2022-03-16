package flags

import (
	"strconv"

	"github.com/spf13/pflag"
)

// todo: add optional flags support:
// https://github.com/dotnet/runtime/blob/main/docs/design/coreclr/botr/xplat-minidump-generation.md

type createdump struct {
	ProcessID int
	Output    string

	flagSet *pflag.FlagSet
}

func NewCreateDump() DotnetTool {
	return &createdump{
		ProcessID: 1,
	}
}

func (cd *createdump) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(cd.BinaryName(), pflag.ExitOnError)
	flagSet.IntVarP(
		&cd.ProcessID,
		"process-id",
		"p",
		cd.ProcessID,
		"The process ID to collect the trace from",
	)
	flagSet.StringVarP(
		&cd.Output,
		"output",
		"o",
		cd.Output,
		"Output file",
	)
	cd.flagSet = flagSet
	return cd.flagSet
}

func (cd *createdump) FormatArgs(args *Args) {
	args.AppendRaw(strconv.Itoa(cd.ProcessID))

	if cd.Output != "" {
		args.Append("name", cd.Output)
	}
}

func (*createdump) SetAction(_ string) {
	// omit action usage here because of create dump is not a dotnet tool (it's a runtime tool)
}

func (cd *createdump) SetOutput(output string) {
	cd.Output = output
}

func (cd *createdump) SetProcessID(processID int) {
	cd.ProcessID = processID
}

func (*createdump) BinaryName() string {
	return "createdump"
}

func (*createdump) ToolName() string {
	return "full-dump"
}
