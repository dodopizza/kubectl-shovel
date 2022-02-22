package watchdog

import (
	"context"
	"errors"
	"time"
)

type Operator struct {
	ping     chan struct{}
	check    func() bool
	deadline time.Duration
	interval time.Duration
}

func NewOperator(check func() bool, deadline, interval time.Duration) *Operator {
	return &Operator{
		check:    check,
		deadline: deadline,
		interval: interval,
	}
}

func (p *Operator) Run(ctx context.Context) error {
	p.ping = make(chan struct{}, 1)

	go p.awaiter(ctx)

	for {
		select {
		case <-p.ping:
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
	defer close(p.ping)

	for {
		select {
		case <-ticker.C:
			if p.check() {
				p.ping <- struct{}{}
			}
		case <-ctx.Done():
			return
		}
	}
}
