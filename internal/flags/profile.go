package flags

import (
	"fmt"
	"strings"
)

type Profile string

var (
	supportedProfiles = []string{
		"cpu-sampling",
		"gc-verbose",
		"gc-collect",
	}
)

func (p *Profile) String() string {
	return string(*p)
}

func (p *Profile) Set(str string) error {
	if strings.TrimSpace(str) == "" {
		return fmt.Errorf("no profile given, must be one of: [%s]",
			strings.Join(supportedProfiles, ", "))
	}
	if !ContainsItemString(supportedProfiles, str) {
		return fmt.Errorf("unsupported profile \"%s\", must be one of: [%s]",
			str, strings.Join(supportedProfiles, ", "))

	}
	*p = Profile(str)
	return nil
}

func (p *Profile) Type() string {
	return "profile"
}

func (p *Profile) Description() string {
	return "A named pre-defined set of provider configurations that allows" +
		"common tracing scenarios to be specified succinctly.\n" +
		"The following profiles are available:\n" +
		strings.Join(supportedProfiles, ", ")
}
