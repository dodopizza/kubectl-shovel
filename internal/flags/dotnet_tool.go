package flags

import (
	"fmt"
	"github.com/spf13/pflag"
	"strconv"
)

type DotnetToolFactory func() DotnetTool

type DotnetToolFlags interface {
	GetFlags() *pflag.FlagSet
	FormatArgs() []string
	SetAction(action string) DotnetToolFlags
	SetOutput(output string) DotnetToolFlags
	SetProcessID(id int) DotnetToolFlags
}

type DotnetTool interface {
	DotnetToolFlags
	BinaryName() string
	ToolName() string
	GetProperties() DotnetToolFlags
}

type DotnetToolProperties struct {
	Action    string
	Output    string
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
	flagSet.StringVarP(
		&dt.Output,
		"output",
		"o",
		dt.Output,
		"Output file",
	)
	return flagSet
}

func (dt *DotnetToolProperties) FormatArgs() []string {
	args := make([]string, 0)

	if dt.Action != "" {
		args = append(args, dt.Action)
	}

	args = append(args, "--process-id", strconv.Itoa(dt.ProcessID))

	if dt.Output != "" {
		args = append(args, "--output", dt.Output)
	}

	return args
}

func (dt *DotnetToolProperties) SetAction(action string) DotnetToolFlags {
	dt.Action = action
	return dt
}

func (dt *DotnetToolProperties) SetOutput(output string) DotnetToolFlags {
	dt.Output = output
	return dt
}

func (dt *DotnetToolProperties) SetProcessID(id int) DotnetToolFlags {
	dt.ProcessID = id
	return dt
}

func FlagToArg(flag pflag.Value) []string {
	return []string{
		fmt.Sprintf("--%s", flag.Type()),
		flag.String(),
	}
}
