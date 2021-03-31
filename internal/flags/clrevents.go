package flags

import (
	"fmt"
	"strings"
)

type CLREvents []string

const (
	CLREventsDelimiter = "+"
)

func (c *CLREvents) String() string {
	return strings.Join(*c, CLREventsDelimiter)
}

func (c *CLREvents) Set(str string) error {
	if strings.TrimSpace(str) == "" {
		return fmt.Errorf("passed empty CLR Events")
	}
	events := strings.Split(str, CLREventsDelimiter)
	*c = events
	return nil
}

func (c *CLREvents) Type() string {
	return "clrevents"
}

func (c *CLREvents) Description() string {
	return fmt.Sprintf(
		"A list of CLR runtime provider keywords to enable separated by \"%s\" signs.\n"+
			"This is a simple mapping that lets you specify event keywords via string aliases rather than their hex values.\n"+
			"More info here: https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-trace#options-1",
		CLREventsDelimiter,
	)
}
