package flags

import (
	"fmt"

	"github.com/spf13/pflag"
)

type DotnetToolFlagSetFactory func() DotnetToolFlagSet

type DotnetToolFlagSet interface {
	GetFlags() *pflag.FlagSet
	GetArgs() []string
}

func FlagToArg(flag pflag.Value) []string {
	return []string{
		fmt.Sprintf("--%s", flag.Type()),
		flag.String(),
	}
}
