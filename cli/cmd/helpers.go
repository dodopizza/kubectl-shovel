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

func newJobName() string {
	return strings.Join(
		[]string{
			pluginName,
			uuid.NewString(),
		},
		"-",
	)
}
