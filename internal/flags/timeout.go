package flags

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Timeout time.Duration

func (t *Timeout) String() string {
	return fmt.Sprintf(
		"%.0f",
		time.Duration(*t).
			Seconds(),
	)
}

func (t *Timeout) Set(str string) error {
	if strings.TrimSpace(str) == "" {
		return fmt.Errorf("empty timeout passed")
	}
	d, err := time.ParseDuration(str)
	if err == nil {
		*t = Timeout(d)
		return t.checkDuration()
	}
	i, err := strconv.Atoi(str)
	if err == nil {
		*t = Timeout(time.Duration(i) * time.Second)
		return t.checkDuration()
	}

	return fmt.Errorf(
		"can't parse duration from %s, provide it as number to define seconds or with units",
		str,
	)
}

func (t *Timeout) checkDuration() error {
	if time.Duration(*t) < (1 * time.Second) {
		return fmt.Errorf("provided duration is to low, minimum 1 second")
	}
	return nil
}

func (t *Timeout) Type() string {
	return "timeout"
}

func (t *Timeout) Description() string {
	return "Give up on collecting the GC dump if it takes longer than this many seconds.\n" +
		"Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".\n" +
		"Will be rounded to seconds. If no unit provided defaults to seconds.\n" +
		"(default 30 sec)"
}
