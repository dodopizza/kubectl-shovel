package flags

import (
	"strconv"

	"github.com/dodopizza/kubectl-shovel/internal/flags/types"

	"github.com/spf13/pflag"
)

type trace struct {
	*DotnetToolSharedOptions

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
	return &trace{
		DotnetToolSharedOptions: NewDotnetToolSharedOptions(),
		BufferSize:              256,
	}
}

func (tr *trace) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(tr.BinaryName(), pflag.ExitOnError)
	flagSet.AddFlagSet(tr.DotnetToolSharedOptions.GetFlags())
	flagSet.IntVar(
		&tr.BufferSize,
		"buffersize",
		tr.BufferSize,
		"Sets the size of the in-memory circular buffer, in megabytes",
	)

	flagSet.Var(
		&tr.CLREventLevel,
		tr.CLREventLevel.Type(),
		tr.CLREventLevel.Description(),
	)

	flagSet.Var(
		&tr.CLREvents,
		tr.CLREvents.Type(),
		tr.CLREvents.Description(),
	)

	tr.Duration = types.Duration(types.DefaultDuration)
	flagSet.Var(
		&tr.Duration,
		tr.Duration.Type(),
		tr.Duration.Description(),
	)

	flagSet.Var(
		&tr.Format,
		tr.Format.Type(),
		tr.Format.Description(),
	)

	flagSet.Var(
		&tr.Profile,
		tr.Profile.Type(),
		tr.Profile.Description(),
	)

	flagSet.Var(
		&tr.Providers,
		tr.Providers.Type(),
		tr.Providers.Description(),
	)

	tr.flagSet = flagSet
	return flagSet
}

func (tr *trace) FormatArgs(args *Args, t FormatArgsType) {
	tr.DotnetToolSharedOptions.FormatArgs(args, t)

	if tr.flagSet.Changed("buffersize") {
		args.Append("buffersize", strconv.Itoa(tr.BufferSize))
	}

	if tr.flagSet.Changed(tr.CLREventLevel.Type()) {
		args.AppendFlag(&tr.CLREventLevel)
	}

	if tr.flagSet.Changed(tr.CLREvents.Type()) {
		args.AppendFlag(&tr.CLREvents)
	}

	args.AppendFlag(&tr.Duration)

	if tr.flagSet.Changed(tr.Format.Type()) {
		args.AppendFlag(&tr.Format)
	}

	if tr.flagSet.Changed(tr.Profile.Type()) {
		args.AppendFlag(&tr.Profile)
	}

	if tr.flagSet.Changed(tr.Providers.Type()) {
		args.AppendFlag(&tr.Providers)
	}
}

func (*trace) BinaryName() string {
	return "dotnet-trace"
}

func (*trace) ToolName() string {
	return "trace"
}

func (*trace) IsLimitedResources() bool {
	return true
}
