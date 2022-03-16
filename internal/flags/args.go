package flags

import (
	"github.com/spf13/pflag"
)

// Args contains cli argument state
type Args struct {
	items []string
}

// NewArgs returns new instance with empty arguments
func NewArgs() *Args {
	return &Args{
		items: make([]string, 0),
	}
}

// Append adds argument --argument with corresponding value to state
func (a *Args) Append(argument, value string) *Args {
	a.items = append(a.items, "--"+argument, value)
	return a
}

// AppendRaw adds string item to state
func (a *Args) AppendRaw(item string) *Args {
	a.items = append(a.items, item)
	return a
}

// AppendKey adds argument --argument without value to state
func (a *Args) AppendKey(argument string) *Args {
	a.items = append(a.items, "--"+argument)
	return a
}

// AppendFlag adds argument pflag.Value.Type() with corresponding value pflag.Value.String() to state
func (a *Args) AppendFlag(flag pflag.Value) *Args {
	return a.Append(flag.Type(), flag.String())
}

// Get returns args state as slice of strings
func (a *Args) Get() []string {
	return a.items
}
