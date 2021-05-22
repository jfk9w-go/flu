package flu

import (
	"context"
	"sync"
)

type WaitGroup struct {
	sync.WaitGroup
}

func (wg *WaitGroup) Go(ctx context.Context, fun func(ctx context.Context)) func() {
	wg.Add(1)
	ctx, cancel := context.WithCancel(ctx)
	go func(ctx context.Context, cancel func()) {
		defer wg.Done()
		defer cancel()
		fun(ctx)
	}(ctx, cancel)
	return cancel
}
