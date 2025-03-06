package flags

import (
	"github.com/spf13/pflag"
)

const (
	FormatArgsTypeBinary = FormatArgsType("binary")
	FormatArgsTypeTool   = FormatArgsType("tool")
)

type DotnetToolFactory func() DotnetTool
type FormatArgsType string

type DotnetTool interface {
	DotnetToolFlagsFormatter
	BinaryName() string
	ToolName() string
	IsPrivileged() bool
	IsLimitedResources() bool
}

type DotnetToolFlagsFormatter interface {
	FormatArgs(args *Args, t FormatArgsType)
	GetFlags() *pflag.FlagSet
	SetOutput(output string)
	SetProcessID(id int)
}
