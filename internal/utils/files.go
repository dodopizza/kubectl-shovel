package utils

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

// MoveFile physically moves file from source path to destination path
// If dest already exists, MoveFile replaces it
func MoveFile(source, dest string) error {
	output, _ := os.Create(dest)
	defer Ignore(output.Close)

	input, _ := os.Open(source)
	_, err := io.Copy(output, input)
	_ = input.Close()
	if err != nil {
		return errors.Wrapf(err, "failed to move from: %s to: %s", source, dest)
	}

	err = os.Remove(source)
	if err != nil {
		return err
	}
	return nil
}

// FileExists returns value indicating that file exists on file system
func FileExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}

	return errors.Is(err, os.ErrNotExist)
}
