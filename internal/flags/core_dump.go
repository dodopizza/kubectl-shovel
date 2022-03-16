package flags

import (
	"strconv"
)

// todo: add optional flags support:
// https://github.com/dotnet/runtime/blob/main/docs/design/coreclr/botr/xplat-minidump-generation.md

type coredump struct {
	*DotnetToolSharedOptions
}

func NewCoreDump() DotnetTool {
	return &coredump{
		DotnetToolSharedOptions: NewDotnetToolSharedOptions(),
	}
}

func (cd *coredump) FormatArgs(args *Args, t FormatArgsType) {
	// preserve same args interface for all available commands
	// but format correct args for binary

	if t == FormatArgsTypeTool {
		cd.DotnetToolSharedOptions.FormatArgs(args, t)
		return
	}

	args.AppendRaw(strconv.Itoa(cd.ProcessID))

	if cd.Output != "" {
		args.Append("name", cd.Output)
	}
}

func (*coredump) BinaryName() string {
	return "createdump"
}

func (*coredump) ToolName() string {
	return "coredump"
}

func (*coredump) IsPrivileged() bool {
	return true
}
