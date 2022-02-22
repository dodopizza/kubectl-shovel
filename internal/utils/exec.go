package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// ExecCommand is wrapper for running external processes
func ExecCommand(executable string, args ...string) error {
	cmd := exec.Command(
		executable,
		args...,
	)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, stdout.String())
	}

	return nil
}

// MoveFile physically moves file from source path to destination path
// If dest already exists, MoveFile replaces it
func MoveFile(source, dest string) error {
	output, _ := os.Create(dest)
	defer output.Close()

	input, _ := os.Open(source)
	_, err := io.Copy(output, input)
	_ = input.Close()
	if err != nil {
		return fmt.Errorf("move failed: %s", err)
	}

	err = os.Remove(source)
	if err != nil {
		return fmt.Errorf("move failed: %s", err)
	}
	return nil
}
