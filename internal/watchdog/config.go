package watchdog

import (
	"time"
)

const (
	deadAfterDuration = 1 * time.Minute
	pingInterval      = 10 * time.Second
	checkInterval     = 5 * time.Second

	pingFile = "/ping.tmp"
)
