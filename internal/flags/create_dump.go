package flags

import (
	"strconv"
)

// todo: add optional flags support:
// https://github.com/dotnet/runtime/blob/main/docs/design/coreclr/botr/xplat-minidump-generation.md

type createdump struct {
	*DotnetToolSharedOptions
}

func NewCreateDump() DotnetTool {
	return &createdump{
		DotnetToolSharedOptions: NewDotnetToolSharedOptions(),
	}
}

func (cd *createdump) FormatArgs(args *Args, t FormatArgsType) {
	// preserve same args for all available commands
	// but format correct args for binary execution

	if t == FormatArgsTypeTool {
		cd.DotnetToolSharedOptions.FormatArgs(args, t)
		return
	}

	args.AppendRaw(strconv.Itoa(cd.ProcessID))

	if cd.Output != "" {
		args.Append("name", cd.Output)
	}
}

func (*createdump) SetAction(_ string) {
	// omit action usage here because of create dump is not a dotnet tool (it's a runtime tool)
}

func (*createdump) BinaryName() string {
	return "createdump"
}

func (*createdump) ToolName() string {
	return "full-dump"
}
