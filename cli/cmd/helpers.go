package cmd

import (
	"strconv"
	"strings"
	"time"
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
			currentTime(),
		},
		"-",
	)
}

func currentTime() string {
	return strconv.Itoa(
		int(time.Now().Unix()),
	)
}
