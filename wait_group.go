package flu

import (
	"context"
	"sync"
)

type WaitGroup struct {
	sync.WaitGroup
}

func (wg *WaitGroup) Go(ctx context.Context, rateLimiter RateLimiter, fun func(ctx context.Context)) func() {
	if rateLimiter == nil {
		rateLimiter = RateUnlimiter
	}

	wg.Add(1)
	ctx, cancel := context.WithCancel(ctx)
	go func(ctx context.Context, cancel func(), rateLimiter RateLimiter) {
		defer func() {
			cancel()
			wg.Done()
		}()

		if err := rateLimiter.Start(ctx); err != nil {
			return
		}

		defer rateLimiter.Complete()
		fun(ctx)
	}(ctx, cancel, rateLimiter)

	return cancel
}
