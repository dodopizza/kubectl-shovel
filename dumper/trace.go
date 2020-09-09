package main

import (
	"strconv"
)

const (
	dotnetTraceBinary = "dotnet-trace"
)

func makeTrace(pid int, output string) error {
	return runCommand(dotnetTraceBinary,
		"collect",
		"--process-id",
		strconv.Itoa(pid),
		"--output",
		output,
		"--duration",
		"00:00:00:10",
	)
}
