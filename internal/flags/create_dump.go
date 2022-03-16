package flags

import (
	"strconv"

	"github.com/spf13/pflag"
)

type CreateDump struct {
	ProcessID int
	Output    string
	flagSet   *pflag.FlagSet
}

func NewCreateDump() DotnetTool {
	return &CreateDump{
		ProcessID: 1,
		Output:    "/tmp/coredump.%p",
	}
}

func (cd *CreateDump) GetFlags() *pflag.FlagSet {
	return cd.flagSet
}

func (cd *CreateDump) FormatArgs(args *Args) {
	// todo: additional flags
	args.
		AppendCommand(strconv.Itoa(cd.ProcessID)).
		Append("name", cd.Output)
}

func (cd *CreateDump) SetAction(_ string) DotnetToolFlags {
	// omit action here because of create dump is not a dotnet tool
	return cd
}

func (cd *CreateDump) SetOutput(output string) DotnetToolFlags {
	cd.Output = output
	return cd
}

func (cd *CreateDump) SetProcessID(processID int) DotnetToolFlags {
	cd.ProcessID = processID
	return cd
}

func (*CreateDump) BinaryName() string {
	// todo: use full path to tool ?
	return "createdump"
}

func (*CreateDump) ToolName() string {
	return "full-dump"
}

func (cd *CreateDump) GetProperties() DotnetToolFlags {
	return cd
}
