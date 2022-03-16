package flags

import (
	"strconv"

	"github.com/dodopizza/kubectl-shovel/internal/flags/types"

	"github.com/spf13/pflag"
)

type DotnetTrace struct {
	*DotnetToolProperties

	BufferSize    int
	CLREventLevel types.CLREventLevel
	CLREvents     types.CLREvents
	Duration      types.Duration
	Format        types.Format
	Profile       types.Profile
	Providers     types.Providers

	flagSet *pflag.FlagSet
}

func NewDotnetTrace() DotnetTool {
	return &DotnetTrace{
		DotnetToolProperties: NewDotnetToolProperties(),
		BufferSize:           256,
	}
}

func (t *DotnetTrace) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(t.BinaryName(), pflag.ExitOnError)
	flagSet.AddFlagSet(t.DotnetToolProperties.GetFlags())
	flagSet.IntVar(
		&t.BufferSize,
		"buffersize",
		t.BufferSize,
		"Sets the size of the in-memory circular buffer, in megabytes",
	)

	flagSet.Var(
		&t.CLREventLevel,
		t.CLREventLevel.Type(),
		t.CLREventLevel.Description(),
	)

	flagSet.Var(
		&t.CLREvents,
		t.CLREvents.Type(),
		t.CLREvents.Description(),
	)

	t.Duration = types.Duration(types.DefaultDuration)
	flagSet.Var(
		&t.Duration,
		t.Duration.Type(),
		t.Duration.Description(),
	)

	flagSet.Var(
		&t.Format,
		t.Format.Type(),
		t.Format.Description(),
	)

	flagSet.Var(
		&t.Profile,
		t.Profile.Type(),
		t.Profile.Description(),
	)

	flagSet.Var(
		&t.Providers,
		t.Providers.Type(),
		t.Providers.Description(),
	)

	t.flagSet = flagSet
	return flagSet
}

func (t *DotnetTrace) FormatArgs(args *Args) {
	args.AppendFrom(t.DotnetToolProperties)

	if t.flagSet.Changed("buffersize") {
		args.Append("buffersize", strconv.Itoa(t.BufferSize))
	}

	if t.flagSet.Changed(t.CLREventLevel.Type()) {
		args.AppendFlag(&t.CLREventLevel)
	}

	if t.flagSet.Changed(t.CLREvents.Type()) {
		args.AppendFlag(&t.CLREvents)
	}

	args.AppendFlag(&t.Duration)

	if t.flagSet.Changed(t.Format.Type()) {
		args.AppendFlag(&t.Format)
	}

	if t.flagSet.Changed(t.Profile.Type()) {
		args.AppendFlag(&t.Profile)
	}

	if t.flagSet.Changed(t.Providers.Type()) {
		args.AppendFlag(&t.Providers)
	}
}

func (*DotnetTrace) BinaryName() string {
	return "dotnet-trace"
}

func (*DotnetTrace) ToolName() string {
	return "trace"
}
