package metrics

import (
	"context"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jfk9w-go/flu"
	"github.com/pkg/errors"
)

type GraphiteMetric interface {
	Reset() (last float64, set bool)
}

type GraphiteCounter AtomicFloat64

func (c *GraphiteCounter) Inc() {
	c.Add(1)
}

func (c *GraphiteCounter) Add(delta float64) {
	(*AtomicFloat64)(c).Add(delta)
}

func (c *GraphiteCounter) Reset() (float64, bool) {
	zero := float64(0)
	value := (*AtomicFloat64)(c).Swap(zero)
	return value, value != zero
}

type GraphiteGauge AtomicFloat64

func (g *GraphiteGauge) Set(value float64) {
	(*AtomicFloat64)(g).Set(value)
}

func (g *GraphiteGauge) Inc() {
	g.Add(1)
}

func (g *GraphiteGauge) Dec() {
	g.Add(-1)
}

func (g *GraphiteGauge) Add(delta float64) {
	(*AtomicFloat64)(g).Add(delta)
}

func (g *GraphiteGauge) Sub(delta float64) {
	g.Add(-delta)
}

func (g *GraphiteGauge) Reset() (float64, bool) {
	return (*AtomicFloat64)(g).Get(), true
}

type GraphiteClient struct {
	address string
	prefix  string
	metrics map[string]GraphiteMetric
	mu      *sync.RWMutex
	cancel  func()
	work    *sync.WaitGroup
}

func NewGraphiteClient(ctx context.Context, address string, interval time.Duration) GraphiteClient {
	ctx, cancel := context.WithCancel(ctx)
	client := GraphiteClient{
		address: address,
		metrics: make(map[string]GraphiteMetric),
		mu:      new(sync.RWMutex),
		cancel:  cancel,
		work:    new(sync.WaitGroup),
	}

	if interval > 0 {
		client.work.Add(1)
		go func() {
			timer := time.NewTimer(interval)
			defer func() {
				timer.Stop()
				if err := client.FlushValues(time.Now()); err != nil {
					log.Printf("Failed to flush Graphite metrics: %s", err)
				}
				client.work.Done()
			}()

			for {
				select {
				case <-ctx.Done():
					return
				case now := <-timer.C:
					if err := client.FlushValues(now); err != nil {
						log.Printf("Failed to flush Graphite metrics: %s", err)
					}
				}
			}
		}()
	}

	return client
}

func (g GraphiteClient) Close() {
	g.cancel()
	g.work.Wait()
}

func (g GraphiteClient) FlushValues(now time.Time) error {
	b := new(strings.Builder)
	nowstr := strconv.FormatInt(now.Unix(), 10)

	g.mu.RLock()
	for key, metric := range g.metrics {
		value, set := metric.Reset()
		if !set {
			continue
		}

		b.WriteString(key)
		b.WriteRune(' ')
		b.WriteString(strconv.FormatFloat(value, 'f', 9, 64))
		b.WriteRune(' ')
		b.WriteString(nowstr)
		b.WriteRune('\n')
	}

	g.mu.RUnlock()
	if b.Len() == 0 {
		return nil
	}

	conn, err := net.Dial("tcp", g.address)
	if err != nil {
		return errors.Wrap(err, "connect")
	}

	err = flu.EncodeTo(&flu.PlainText{b.String()}, flu.IO{W: conn})
	if err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func (g GraphiteClient) WithPrefix(prefix string) Client {
	if g.prefix != "" {
		g.prefix += "."
	}

	g.prefix += prefix
	return g
}

func (g GraphiteClient) Counter(name string, labels Labels) Counter {
	key := g.makeKey(name, labels)

	g.mu.RLock()
	entry, ok := g.metrics[key]
	g.mu.RUnlock()

	if !ok {
		g.mu.Lock()
		entry, ok = g.metrics[key]
		if !ok {
			entry = new(GraphiteCounter)
			g.metrics[key] = entry
		}
		g.mu.Unlock()
	}

	return entry.(Counter)
}

func (g GraphiteClient) Gauge(name string, labels Labels) Gauge {
	key := g.makeKey(name, labels)

	g.mu.RLock()
	entry, ok := g.metrics[key]
	g.mu.RUnlock()

	if !ok {
		g.mu.Lock()
		entry, ok = g.metrics[key]
		if !ok {
			entry = new(GraphiteGauge)
			g.metrics[key] = entry
		}
		g.mu.Unlock()
	}

	return entry.(Gauge)
}

func (g GraphiteClient) makeKey(name string, labels Labels) string {
	prefix := g.prefix
	if prefix != "" {
		prefix += "."
	}

	values := labels.Path(".", "_")
	prefix += values
	if values != "" {
		prefix += "."
	}

	return prefix + name
}
