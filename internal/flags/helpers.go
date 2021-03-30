package flags

import (
	"fmt"

	"github.com/spf13/pflag"
)

func ContainsItemString(strs []string, item string) bool {
	for _, str := range strs {
		if str == item {
			return true
		}
	}
	return false
}

func FlagToArg(flag pflag.Value) []string {
	return []string{
		fmt.Sprintf("--%s", flag.Type()),
		flag.String(),
	}
}
