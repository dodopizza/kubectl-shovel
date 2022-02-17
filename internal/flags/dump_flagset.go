package flags

import (
	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
	"github.com/spf13/pflag"
)

type DumpFlagSet struct {
	Diagnostics bool
	Type        types.DumpType
	dt          *DotnetToolsFlagSet
	flagSet     *pflag.FlagSet
}

func NewDumpFlagSet() DotnetToolFlagSet {
	return &DumpFlagSet{
		Diagnostics: false,
		Type:        types.DumpTypeFull,
		dt:          NewDotnetToolsFlagSet(),
	}
}

func (dump *DumpFlagSet) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("dotnet-dump", pflag.ExitOnError)
	flagSet.AddFlagSet(dump.dt.GetFlags())
	flagSet.BoolVar(
		&dump.Diagnostics,
		"diag",
		dump.Diagnostics,
		"Enable dump collection diagnostic logging",
	)
	flagSet.Var(
		&dump.Type,
		dump.Type.Type(),
		dump.Type.Description(),
	)

	dump.flagSet = flagSet
	return flagSet
}

func (dump *DumpFlagSet) GetArgs() []string {
	args := dump.dt.GetArgs()
	if dump.flagSet.Changed("diag") {
		args = append(args, "--diag")
	}
	args = append(args, "--type", dump.Type.String())
	return args
}
