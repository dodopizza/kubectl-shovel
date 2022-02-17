package flags

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
	"strconv"

	"github.com/spf13/pflag"
)

type TraceFlagSet struct {
	BufferSize    int
	CLREventLevel types.CLREventLevel
	CLREvents     types.CLREvents
	Duration      types.Duration
	Format        types.Format
	Profile       types.Profile
	Providers     types.Providers

	dt *DotnetToolsFlagSet

	flagSet *pflag.FlagSet
}

func NewTraceFlagSet() FlagSetContainer {
	return &TraceFlagSet{
		BufferSize: 256,
		dt:         NewDotnetToolsFlagSet(),
	}
}

func (trace *TraceFlagSet) Parse() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("dotnet-trace", pflag.ExitOnError)
	flagSet.AddFlagSet(trace.dt.Parse())
	flagSet.IntVar(
		&trace.BufferSize,
		"buffersize",
		trace.BufferSize,
		"Sets the size of the in-memory circular buffer, in megabytes",
	)

	flagSet.Var(
		&trace.CLREventLevel,
		trace.CLREventLevel.Type(),
		trace.CLREventLevel.Description(),
	)

	flagSet.Var(
		&trace.CLREvents,
		trace.CLREvents.Type(),
		trace.CLREvents.Description(),
	)

	trace.Duration = types.Duration(types.DefaultDuration)
	flagSet.Var(
		&trace.Duration,
		trace.Duration.Type(),
		trace.Duration.Description(),
	)

	flagSet.Var(
		&trace.Format,
		trace.Format.Type(),
		trace.Format.Description(),
	)

	flagSet.Var(
		&trace.Profile,
		trace.Profile.Type(),
		trace.Profile.Description(),
	)

	flagSet.Var(
		&trace.Providers,
		trace.Providers.Type(),
		trace.Providers.Description(),
	)

	trace.flagSet = flagSet
	return flagSet
}

func (trace *TraceFlagSet) Args() []string {
	args := trace.dt.Args()

	if trace.flagSet.Changed("buffersize") {
		args = append(
			args,
			[]string{
				"--buffersize", strconv.Itoa(trace.BufferSize),
			}...,
		)
	}

	if trace.flagSet.Changed(trace.CLREventLevel.Type()) {
		args = append(
			args,
			FlagToArg(&trace.CLREventLevel)...,
		)
	}

	if trace.flagSet.Changed(trace.CLREvents.Type()) {
		args = append(
			args,
			FlagToArg(&trace.CLREvents)...,
		)
	}

	args = append(
		args,
		FlagToArg(&trace.Duration)...,
	)

	if trace.flagSet.Changed(trace.Format.Type()) {
		args = append(
			args,
			FlagToArg(&trace.Format)...,
		)
	}

	if trace.flagSet.Changed(trace.Profile.Type()) {
		args = append(
			args,
			FlagToArg(&trace.Profile)...,
		)
	}

	if trace.flagSet.Changed(trace.Providers.Type()) {
		args = append(
			args,
			FlagToArg(&trace.Providers)...,
		)
	}

	return args
}
