package watchdog

import (
	"time"
)

const (
	Deadline      = 1 * time.Minute
	PingInterval  = 10 * time.Second
	CheckInterval = 5 * time.Second

	pingFile = "/ping.tmp"
)
