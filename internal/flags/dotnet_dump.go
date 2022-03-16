package flags

import (
	"github.com/spf13/pflag"

	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
)

type DotnetDump struct {
	*DotnetToolProperties

	Diagnostics bool
	Type        types.DumpType

	flagSet *pflag.FlagSet
}

func NewDotnetDump() DotnetTool {
	return &DotnetDump{
		DotnetToolProperties: NewDotnetToolProperties(),
		Diagnostics:          false,
		Type:                 types.DumpTypeFull,
	}
}

func (d *DotnetDump) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(d.BinaryName(), pflag.ExitOnError)
	flagSet.AddFlagSet(d.DotnetToolProperties.GetFlags())
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

func (d *DotnetDump) FormatArgs(args *Args) {
	args.AppendFrom(d.DotnetToolProperties)
	if d.flagSet.Changed("diag") {
		args.AppendKey("diag")
	}
	args.Append("type", d.Type.String())
}

func (*DotnetDump) BinaryName() string {
	return "dotnet-dump"
}

func (*DotnetDump) ToolName() string {
	return "dump"
}

func (d *DotnetDump) GetProperties() DotnetToolFlags {
	return d.DotnetToolProperties
}
