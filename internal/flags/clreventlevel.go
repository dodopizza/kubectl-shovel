package flags

import (
	"fmt"
	"strings"
)

type CLREventLevel string

var supportedCLREventLevels = []string{
	"logalways",
	"critical",
	"error",
	"warning",
	"informational",
	"verbose",
}

func (c *CLREventLevel) String() string {
	return string(*c)
}

func (c *CLREventLevel) Set(str string) error {
	if strings.TrimSpace(str) == "" {
		return fmt.Errorf("no CLR Event Level given, must be one of: [%s]",
			strings.Join(supportedCLREventLevels, ", "))
	}
	if !ContainsItemString(supportedCLREventLevels, str) {
		return fmt.Errorf("unsupported CLR Event Level \"%s\", must be one of: [%s]",
			str, strings.Join(supportedCLREventLevels, ", "))

	}
	*c = CLREventLevel(str)
	return nil
}

func (c *CLREventLevel) Type() string {
	return "clreventlevel"
}

func (c *CLREventLevel) Description() string {
	return "Verbosity of CLR events to be emitted. Supported levels:\n" +
		strings.Join(supportedCLREventLevels, ", ")
}
