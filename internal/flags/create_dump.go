package flags

import (
	"strconv"

	"github.com/spf13/pflag"
)

// todo: add optional flags support:
// https://github.com/dotnet/runtime/blob/main/docs/design/coreclr/botr/xplat-minidump-generation.md

type CreateDump struct {
	ProcessID int
	Output    string

	flagSet *pflag.FlagSet
}

func NewCreateDump() DotnetTool {
	return &CreateDump{
		ProcessID: 1,
	}
}

func (cd *CreateDump) GetFlags() *pflag.FlagSet {
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

func (cd *CreateDump) FormatArgs(args *Args) {
	args.AppendRaw(strconv.Itoa(cd.ProcessID))

	if cd.Output != "" {
		args.Append("name", cd.Output)
	}
}

func (*CreateDump) SetAction(_ string) {
	// omit action usage here because of create dump is not a dotnet tool (it's a runtime tool)
}

func (cd *CreateDump) SetOutput(output string) {
	cd.Output = output
}

func (cd *CreateDump) SetProcessID(processID int) {
	cd.ProcessID = processID
}

func (*CreateDump) BinaryName() string {
	return "createdump"
}

func (*CreateDump) ToolName() string {
	return "full-dump"
}
