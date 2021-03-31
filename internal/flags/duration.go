package flags

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Duration time.Duration

var (
	defaultDuration = 10 * time.Second
	day             = 24 * time.Hour
	unitPosition    = map[int]time.Duration{
		0: 24 * time.Hour,
		1: time.Hour,
		2: time.Minute,
		3: time.Second,
	}
)

func (d *Duration) String() string {
	m := time.Duration(*d)
	days := m / day
	m -= days * day
	hours := m / time.Hour
	m -= hours * time.Hour
	minutes := m / time.Minute
	m -= minutes * time.Minute
	seconds := m / time.Second
	return fmt.Sprintf(
		"%02d:%02d:%02d:%02d",
		days, hours, minutes, seconds,
	)
}

func (d *Duration) Set(str string) error {
	if strings.TrimSpace(str) == "" {
		return fmt.Errorf("empty duration passed")
	}
	m, err := time.ParseDuration(str)
	if err == nil {
		*d = Duration(m)
		return d.checkDuration()
	}

	split := strings.Split(str, ":")
	if len(split) != 4 {
		return fmt.Errorf(
			"wrong duration format provided \"%s\", should be \"dd:hh:mm:ss\"",
			str,
		)
	}
	for i := range split {
		n, err := strconv.Atoi(split[i])
		if err != nil {
			return fmt.Errorf(
				"not a number value provided for duration: \"%s\"",
				str,
			)
		}
		m += time.Duration(n) * unitPosition[i]
	}
	*d = Duration(m)
	return d.checkDuration()
}

func (d *Duration) checkDuration() error {
	if time.Duration(*d) < (1 * time.Second) {
		return fmt.Errorf("provided duration is to low, minimum 1 second")
	}
	return nil
}

func (d *Duration) Type() string {
	return "duration"
}

func (d *Duration) Description() string {
	return "Trace for the given timespan and then automatically stop the trace." +
		fmt.Sprintf(
			"Provided in the form of dd:hh:mm:ss or "+
				"corresponding time unit representation (e.g. 1s, 2m, 3h) (default %s)",
			defaultDuration,
		)
}
