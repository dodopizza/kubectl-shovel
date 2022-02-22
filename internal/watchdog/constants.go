package watchdog

import (
	"time"
)

const (
	deadline      = 1 * time.Minute
	pingInterval  = 10 * time.Second
	checkInterval = 5 * time.Second

	pingFile = "/ping.tmp"
)
