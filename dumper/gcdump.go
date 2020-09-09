package main

import (
	"strconv"
)

const (
	dotnetGCDumpBinary = "dotnet-gcdump"
)

func makeGCDump(pid int, output string) error {
	return runCommand(dotnetGCDumpBinary,
		"collect",
		"--process-id",
		strconv.Itoa(pid),
		"--output",
		output,
	)
}
