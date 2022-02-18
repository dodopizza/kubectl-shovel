package flags

import (
	"fmt"
	"strconv"

	"github.com/spf13/pflag"
)

type DotnetToolFactory func() DotnetTool

type DotnetTool interface {
	GetFlags() *pflag.FlagSet
	GetArgs() []string
	BinaryName() string
	ToolName() string
}

type DotnetToolShared struct {
	ProcessID int
}

func NewDotnetToolShared() *DotnetToolShared {
	return &DotnetToolShared{
		ProcessID: 1,
	}
}

func (dt *DotnetToolShared) GetFlags() *pflag.FlagSet {
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

func (dt *DotnetToolShared) GetArgs() []string {
	return []string{
		"collect", "--process-id", strconv.Itoa(dt.ProcessID),
	}
}

func FlagToArg(flag pflag.Value) []string {
	return []string{
		fmt.Sprintf("--%s", flag.Type()),
		flag.String(),
	}
}
