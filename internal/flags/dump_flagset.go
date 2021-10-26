package flags

import (
	"github.com/spf13/pflag"
)

type DumpFlagSet struct {
	Diagnostics bool
	Type        DumpType
	dt          *DotnetToolsFlagSet
	flagSet     *pflag.FlagSet
}

func NewDumpFlagSet() *DumpFlagSet {
	return &DumpFlagSet{
		Diagnostics: false,
		Type:        DumpTypeFull,
		dt:          NewDotnetToolsFlagSet(),
	}
}

func (dump *DumpFlagSet) Parse() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("dotnet-dump", pflag.ExitOnError)
	flagSet.AddFlagSet(dump.dt.Parse())
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

func (dump *DumpFlagSet) Args() []string {
	args := dump.dt.Args()
	if dump.flagSet.Changed("diag") {
		args = append(args, "--diag")
	}
	args = append(args, "--type", dump.Type.String())
	return args
}
