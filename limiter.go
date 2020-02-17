package flu

import (
	"context"
	"time"
)

type RateLimiter interface {
	Start(ctx context.Context) error
	Complete()
}

type concurrencyRateLimiter chan bool

func ConcurrencyRateLimiter(concurrency int) RateLimiter {
	if concurrency < 1 {
		return rateUnlimiter{}
	} else {
		r := make(concurrencyRateLimiter, concurrency)
		for i := 0; i < concurrency; i++ {
			r.Complete()
		}

		return r
	}
}

func (lim concurrencyRateLimiter) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-lim:
		return nil
	}
}

func (lim concurrencyRateLimiter) Complete() {
	lim <- true
}

type intervalRateLimiter struct {
	event    chan time.Time
	interval time.Duration
}

func IntervalRateLimiter(interval time.Duration) RateLimiter {
	if interval <= 0 {
		return rateUnlimiter{}
	} else {
		event := make(chan time.Time, 1)
		event <- time.Unix(0, 0)
		return intervalRateLimiter{event, interval}
	}
}

func (lim intervalRateLimiter) Start(ctx context.Context) error {
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

func (lim intervalRateLimiter) Complete() {
	lim.event <- time.Now()
}

var RateUnlimiter RateLimiter = rateUnlimiter{}

type rateUnlimiter struct{}

func (rateUnlimiter) Start(ctx context.Context) error {
	return nil
}

func (rateUnlimiter) Complete() {

}
