package flags

import (
	"fmt"

	"github.com/spf13/pflag"
)

// todo: better naming

type FlagSetContainerFactory func() FlagSetContainer

type FlagSetContainer interface {
	GetFlags() *pflag.FlagSet
	GetArgs() []string
}

func FlagToArg(flag pflag.Value) []string {
	return []string{
		fmt.Sprintf("--%s", flag.Type()),
		flag.String(),
	}
}
