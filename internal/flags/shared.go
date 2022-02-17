package flags

import (
	"fmt"

	"github.com/spf13/pflag"
)

// todo: better naming

type FlagSetContainerFactory func() FlagSetContainer

type FlagSetContainer interface {
	Parse() *pflag.FlagSet
	Args() []string
}

func FlagToArg(flag pflag.Value) []string {
	return []string{
		fmt.Sprintf("--%s", flag.Type()),
		flag.String(),
	}
}
