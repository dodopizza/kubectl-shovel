package flags

import (
	"fmt"
	"strconv"
	"strings"
)

type Duration string

var defaultDuration = "00:00:00:10"

func (d *Duration) String() string {
	return string(*d)
}

func (d *Duration) Set(str string) error {
	split := strings.Split(str, ":")
	if len(split) != 4 {
		return fmt.Errorf(
			"wrong duration format provided \"%s\", should be \"dd:hh:mm:ss\"",
			str,
		)
	}
	for _, s := range split {
		if _, err := strconv.Atoi(s); err != nil {
			return fmt.Errorf(
				"wrong duration format provided \"%s\", should be \"dd:hh:mm:ss\"",
				str,
			)
		}
	}
	*d = Duration(str)
	return nil
}

func (d *Duration) Type() string {
	return "duration"
}

func (d *Duration) Description() string {
	return "Trace for the given timespan and then automatically stop the trace. Provided in the form of dd:hh:mm:ss"
}
