package types

import (
	"fmt"
	"strings"
)

type CoreDumpType string

func (c *CoreDumpType) String() string {
	return string(*c)
}

func (c *CoreDumpType) Set(str string) error {
	if strings.TrimSpace(str) == "" {
		return fmt.Errorf("no Core Dump Type given, must be one of: [%s]",
			strings.Join(supportedDumpTypes, ", "))
	}
	if !ContainsItemString(supportedDumpTypes, str) {
		return fmt.Errorf("unsupported Core Dump Type \"%s\", must be one of: [%s]",
			str, strings.Join(supportedCLREventLevels, ", "))
	}
	*c = CoreDumpType(str)
	return nil
}

func (*CoreDumpType) Type() string {
	return "type"
}

func (c *CoreDumpType) Value() string {
	switch *c {
	case DumpTypeFull:
		return "full"
	case DumpTypeHeap:
		return "withheap"
	case DumpTypeMini:
		return "normal"
	case DumpTypeTriage:
		return "triage"
	default:
		return ""
	}
}

func (*CoreDumpType) Description() string {
	return "The kinds of information that are collected from process. Supported types:\n" +
		strings.Join(supportedDumpTypes, ", ") + "\n" +
		"Full - A full core dump with all process memory\n" +
		"Heap - A process dump with heap\n" +
		"Mini - A normal minidump with some process information\n" +
		"Triage - A small dump containing minimal information"
}
