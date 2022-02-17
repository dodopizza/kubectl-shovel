package flags

import (
	"fmt"

	"github.com/spf13/pflag"
)

func FlagToArg(flag pflag.Value) []string {
	return []string{
		fmt.Sprintf("--%s", flag.Type()),
		flag.String(),
	}
}
