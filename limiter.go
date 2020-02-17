package flu

import (
	"time"
)

type Limiter interface {
	Start()
	Complete()
}

type concurrencyLimiter chan struct{}

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

func (lim concurrencyLimiter) Start() {
	<-lim
}

var unit struct{}

func (lim concurrencyLimiter) Complete() {
	lim <- unit
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

func (lim intervalLimiter) Start() {
	prev := <-lim.event
	time.Sleep(lim.interval - time.Now().Sub(prev))
}

func (lim intervalLimiter) Complete() {
	lim.event <- time.Now()
}

var Unlimiter Limiter = unlimiter{}

type unlimiter struct{}

func (unlimiter) Start() {

}

func (unlimiter) Complete() {

}
