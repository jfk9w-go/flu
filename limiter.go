package flu

import (
	"context"
	"time"
)

type Limiter interface {
	Start(ctx context.Context) error
	Complete()
}

type concurrencyLimiter chan bool

func ConcurrencyLimiter(concurrency int) Limiter {
	if concurrency < 1 {
		return unlimiter{}
	} else {
		r := make(concurrencyLimiter, concurrency)
		for i := 0; i < concurrency; i++ {
			r.Complete()
		}

		return r
	}
}

func (lim concurrencyLimiter) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-lim:
		return nil
	}
}

func (lim concurrencyLimiter) Complete() {
	lim <- true
}

type intervalLimiter struct {
	event    chan time.Time
	interval time.Duration
}

func IntervalLimiter(interval time.Duration) Limiter {
	if interval <= 0 {
		return unlimiter{}
	} else {
		event := make(chan time.Time, 1)
		event <- time.Unix(0, 0)
		return intervalLimiter{event, interval}
	}
}

func (lim intervalLimiter) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case prevRun := <-lim.event:
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(lim.interval - time.Now().Sub(prevRun)):
			return nil
		}
	}
}

func (lim intervalLimiter) Complete() {
	lim.event <- time.Now()
}

var Unlimiter Limiter = unlimiter{}

type unlimiter struct{}

func (unlimiter) Start(ctx context.Context) error {
	return nil
}

func (unlimiter) Complete() {

}
