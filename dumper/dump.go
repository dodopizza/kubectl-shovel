package main

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
)

const (
	dotnetGCDumpBinary = "dotnet-gcdump"
)

func makeGcDump(pid int, dumpOutput string) error {
	cmd := exec.Command(
		dotnetGCDumpBinary,
		"collect",
		"-p",
		strconv.Itoa(pid),
		"-o",
		dumpOutput,
	)

	var gcDumpStdOut bytes.Buffer
	cmd.Stdout = &gcDumpStdOut

	err := cmd.Run()

	log.Printf(gcDumpStdOut.String())

	if err != nil {
		return err
	}

	return nil
}
