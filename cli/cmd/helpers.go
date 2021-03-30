package cmd

import (
	"strconv"
	"strings"
	"time"

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

func currentTime() string {
	return strconv.Itoa(
		int(time.Now().Unix()),
	)
}
