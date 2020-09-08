package cmd

import (
	"strings"

	"github.com/google/uuid"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func newUUID() (string, error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	id := string(u.String())

	return id, nil
}

func newJobName() (string, error) {
	id, err := newUUID()
	if err != nil {
		return "", err
	}

	return strings.Join(
		[]string{
			pluginName,
			id,
		},
		"-",
	), nil
}
