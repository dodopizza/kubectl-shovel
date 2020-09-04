package main

import (
	"os/exec"
	"bytes"
	"log"
)

const (
	dotnetGCDumpBinary = "dotnet-gcdump"
)

func makeGcDump(pid, dumpOutput string) (error) {
	cmd := exec.Command(dotnetGCDumpBinary, "collect", "-p", pid, "-o", dumpOutput)

	var gcDumpStdOut bytes.Buffer
	cmd.Stdout = &gcDumpStdOut

	err := cmd.Run()

	log.Printf(gcDumpStdOut.String())

	if err != nil {
		return err
	}

	return nil
}