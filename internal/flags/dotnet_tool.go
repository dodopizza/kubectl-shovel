package flags

import (
	"fmt"
	"strconv"

	"github.com/spf13/pflag"
)

type DotnetToolFactory func() DotnetTool

type DotnetToolFlags interface {
	GetFlags() *pflag.FlagSet
	GetArgs() []string
}

type DotnetTool interface {
	DotnetToolFlags
	BinaryName() string
	ToolName() string
	GetProperties() *DotnetToolProperties
}

type DotnetToolProperties struct {
	ProcessID int
}

func NewDotnetToolProperties() *DotnetToolProperties {
	return &DotnetToolProperties{
		ProcessID: 1,
	}
}

func (dt *DotnetToolProperties) GetFlags() *pflag.FlagSet {
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

func (dt *DotnetToolProperties) GetArgs() []string {
	return []string{
		"--process-id", strconv.Itoa(dt.ProcessID),
	}
}

func FlagToArg(flag pflag.Value) []string {
	return []string{
		fmt.Sprintf("--%s", flag.Type()),
		flag.String(),
	}
}
