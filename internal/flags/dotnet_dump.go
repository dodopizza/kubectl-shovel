package flags

import (
	"github.com/spf13/pflag"

	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
)

type dump struct {
	*DotnetToolSharedOptions

	Diagnostics bool
	Type        types.DumpType

	flagSet *pflag.FlagSet
}

func NewDotnetDump() DotnetTool {
	return &dump{
		DotnetToolSharedOptions: NewDotnetToolProperties(),
		Diagnostics:             false,
		Type:                    types.DumpTypeFull,
	}
}

func (d *dump) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(d.BinaryName(), pflag.ExitOnError)
	flagSet.AddFlagSet(d.DotnetToolSharedOptions.GetFlags())
	flagSet.BoolVar(
		&d.Diagnostics,
		"diag",
		d.Diagnostics,
		"Enable dump collection diagnostic logging",
	)
	flagSet.Var(
		&d.Type,
		d.Type.Type(),
		d.Type.Description(),
	)

	d.flagSet = flagSet
	return flagSet
}

func (d *dump) FormatArgs(args *Args) {
	d.DotnetToolSharedOptions.FormatArgs(args)
	if d.flagSet.Changed("diag") {
		args.AppendKey("diag")
	}
	args.Append("type", d.Type.String())
}

func (*dump) BinaryName() string {
	return "dotnet-dump"
}

func (*dump) ToolName() string {
	return "dump"
}
