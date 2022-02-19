package flags

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
	"strconv"

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

func (t *DotnetTrace) FormatArgs() []string {
	args := t.DotnetToolProperties.FormatArgs()

	if t.flagSet.Changed("buffersize") {
		args = append(
			args,
			[]string{
				"--buffersize", strconv.Itoa(t.BufferSize),
			}...,
		)
	}

	if t.flagSet.Changed(t.CLREventLevel.Type()) {
		args = append(
			args,
			FlagToArg(&t.CLREventLevel)...,
		)
	}

	if t.flagSet.Changed(t.CLREvents.Type()) {
		args = append(
			args,
			FlagToArg(&t.CLREvents)...,
		)
	}

	args = append(
		args,
		FlagToArg(&t.Duration)...,
	)

	if t.flagSet.Changed(t.Format.Type()) {
		args = append(
			args,
			FlagToArg(&t.Format)...,
		)
	}

	if t.flagSet.Changed(t.Profile.Type()) {
		args = append(
			args,
			FlagToArg(&t.Profile)...,
		)
	}

	if t.flagSet.Changed(t.Providers.Type()) {
		args = append(
			args,
			FlagToArg(&t.Providers)...,
		)
	}

	return args
}

func (t *DotnetTrace) BinaryName() string {
	return "dotnet-trace"
}

func (t *DotnetTrace) ToolName() string {
	return "trace"
}

func (t *DotnetTrace) GetProperties() DotnetToolFlags {
	return t.DotnetToolProperties
}
