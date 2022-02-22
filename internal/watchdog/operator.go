package watchdog

import (
	"context"
	"errors"
	"time"
)

// Operator represents operator state
type Operator struct {
	check    func() bool
	deadline time.Duration
	interval time.Duration
	signal   chan struct{}
}

// NewOperator returns new operator with specified check function, deadline and interval durations
func NewOperator(check func() bool, deadline, interval time.Duration) *Operator {
	return &Operator{
		check:    check,
		deadline: deadline,
		interval: interval,
	}
}

// Run starts operator
func (p *Operator) Run(ctx context.Context) error {
	p.signal = make(chan struct{}, 1)

	go p.awaiter(ctx)

	for {
		select {
		case <-p.signal:
		case <-ctx.Done():
			return nil
		case <-time.After(p.deadline):
			return errors.New("there were no signals from operator for a long time")
		}
	}
}

func (p *Operator) awaiter(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	defer close(p.signal)

	for {
		select {
		case <-ticker.C:
			if p.check() {
				p.signal <- struct{}{}
			}
		case <-ctx.Done():
			return
		}
	}
}
