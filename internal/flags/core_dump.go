package flags

import (
	"strconv"

	"github.com/spf13/pflag"

	"github.com/dodopizza/kubectl-shovel/internal/flags/types"
)

type coredump struct {
	*DotnetToolSharedOptions

	Type types.CoreDumpType

	flagSet *pflag.FlagSet
}

func NewCoreDump() DotnetTool {
	return &coredump{
		DotnetToolSharedOptions: NewDotnetToolSharedOptions(),
		Type:                    types.DumpTypeFull,
	}
}

func (cd *coredump) GetFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(cd.BinaryName(), pflag.ExitOnError)
	flagSet.AddFlagSet(cd.DotnetToolSharedOptions.GetFlags())
	flagSet.Var(
		&cd.Type,
		cd.Type.Type(),
		cd.Type.Description(),
	)

	cd.flagSet = flagSet
	return flagSet
}

func (cd *coredump) FormatArgs(args *Args, t FormatArgsType) {
	// preserve same args interface for all available commands
	// but format correct args for binary

	if t == FormatArgsTypeTool {
		cd.DotnetToolSharedOptions.FormatArgs(args, t)
		args.Append("type", cd.Type.String())
		return
	}

	args.AppendRaw(strconv.Itoa(cd.ProcessID))
	args.AppendKey(cd.Type.Value())

	if cd.Output != "" {
		args.Append("name", cd.Output)
	}
}

func (*coredump) BinaryName() string {
	return "createdump"
}

func (*coredump) ToolName() string {
	return "coredump"
}

func (*coredump) IsPrivileged() bool {
	return true
}
