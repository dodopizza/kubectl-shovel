package types

import (
	"fmt"
	"strings"
)

type DumpType string

const (
	DumpTypeFull   = "Full"
	DumpTypeHeap   = "Heap"
	DumpTypeMini   = "Mini"
	DumpTypeTriage = "Triage"
)

var supportedDumpTypes = []string{
	DumpTypeFull,
	DumpTypeHeap,
	DumpTypeMini,
	DumpTypeTriage,
}

func (dt *DumpType) String() string {
	return string(*dt)
}

func (dt *DumpType) Set(str string) error {
	if strings.TrimSpace(str) == "" {
		return fmt.Errorf("no Dump Type given, must be one of: [%s]",
			strings.Join(supportedDumpTypes, ", "))
	}
	if !ContainsItemString(supportedDumpTypes, str) {
		return fmt.Errorf("unsupported Dump Type \"%s\", must be one of: [%s]",
			str, strings.Join(supportedCLREventLevels, ", "))
	}
	*dt = DumpType(str)
	return nil
}

func (*DumpType) Type() string {
	return "type"
}

// revive:disable:line-length-limit This is an extended description
func (*DumpType) Description() string {
	return "The kinds of information that are collected from process. Supported types:\n" +
		strings.Join(supportedDumpTypes, ", ") + "\n" +
		"Full - The largest dump containing all memory including the module images\n" +
		"Heap - A large and relatively comprehensive dump containing module lists, thread lists, all stacks, exception information and all memory except for mapped images\n" +
		"Mini - A small dump containing module lists, thread lists, exception information and all stacks\n" +
		"Triage - A small dump containing minimal information"
}

// revive:enable:line-length-limit
