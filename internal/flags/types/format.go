package types

import (
	"fmt"
	"strings"
)

type Format string

var (
	defaultFormat    = "NetTrace"
	supportedFormats = []string{
		defaultFormat,
		"Chromium",
		"Speedscope",
	}
)

func (f *Format) String() string {
	return string(*f)
}

func (f *Format) Set(str string) error {
	if strings.TrimSpace(str) == "" {
		return fmt.Errorf("no format given, must be one of: [%s]",
			strings.Join(supportedFormats, ", "))
	}
	if !ContainsItemString(supportedFormats, str) {
		return fmt.Errorf("unsupported format \"%s\", must be one of: [%s]",
			str, strings.Join(supportedFormats, ", "))
	}
	*f = Format(str)
	return nil
}

func (*Format) Type() string {
	return "format"
}

func (*Format) Description() string {
	return "Sets the output format for the trace file conversion. Supported formats:\n" +
		strings.Join(supportedFormats, ", ") +
		fmt.Sprintf(" (default \"%s\")", defaultFormat)
}
