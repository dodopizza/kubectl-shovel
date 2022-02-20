package flags

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
	"github.com/spf13/pflag"
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
	flagSet := pflag.NewFlagSet("dotnet-dump", pflag.ExitOnError)
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

func (d *DotnetDump) FormatArgs() []string {
	args := d.DotnetToolProperties.FormatArgs()
	if d.flagSet.Changed("diag") {
		args = append(args, "--diag")
	}
	args = append(args, "--type", d.Type.String())
	return args
}

func (d *DotnetDump) BinaryName() string {
	return "dotnet-dump"
}

func (d *DotnetDump) ToolName() string {
	return "dump"
}

func (d *DotnetDump) GetProperties() DotnetToolFlags {
	return d.DotnetToolProperties
}
