package flags

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
	"github.com/spf13/pflag"
)

type DotnetDump struct {
	Diagnostics bool
	Type        types.DumpType
	dt          *DotnetToolShared

	flagSet *pflag.FlagSet
}

func NewDotnetDump() DotnetTool {
	return &DotnetDump{
		Diagnostics: false,
		Type:        types.DumpTypeFull,
		dt:          NewDotnetToolShared(),
	}
}

func (d *DotnetDump) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("dotnet-dump", pflag.ExitOnError)
	flagSet.AddFlagSet(d.dt.GetFlags())
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

func (d *DotnetDump) GetArgs() []string {
	args := d.dt.GetArgs()
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
