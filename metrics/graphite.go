package metrics

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jfk9w-go/flu"
	"github.com/pkg/errors"
)

var GraphiteTimeout = 1 * time.Minute

type GraphiteMetric interface {
	Write(b *strings.Builder, now string, key string)
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

func (c *GraphiteCounter) Write(b *strings.Builder, now string, key string) {
	value, set := c.Reset()
	if !set {
		return
	}

	b.WriteString(key)
	b.WriteRune(' ')
	b.WriteString(strconv.FormatFloat(value, 'f', 9, 64))
	b.WriteRune(' ')
	b.WriteString(now)
	b.WriteRune('\n')
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

func (g *GraphiteGauge) Write(b *strings.Builder, now string, key string) {
	value, set := g.Reset()
	if !set {
		return
	}

	b.WriteString(key)
	b.WriteRune(' ')
	b.WriteString(strconv.FormatFloat(value, 'f', 9, 64))
	b.WriteRune(' ')
	b.WriteString(now)
	b.WriteRune('\n')
}

type GraphiteHistogram struct {
	buckets  []float64
	counters []*GraphiteCounter
	hbf      string
}

func (h GraphiteHistogram) Observe(value float64) {
	idx := len(h.buckets)
	for i, upper := range h.buckets {
		if value < upper {
			idx = i
			break
		}
	}

	h.counters[idx].Inc()
}

func (h GraphiteHistogram) Write(b *strings.Builder, now string, key string) {
	for i, counter := range h.counters {
		bucket := "inf"
		if h.buckets[i] != math.MaxFloat64 {
			bucket = fmt.Sprintf(h.hbf, h.buckets[i])
		}

		counter.Write(b, now, key+"."+strings.Replace(bucket, ".", "_", 1))
	}
}

type GraphiteClient struct {
	HistogramBucketFormat string

	address string
	prefix  string
	metrics map[string]GraphiteMetric
	cancel  func()

	mu   *flu.RWMutex
	work *flu.WaitGroup
}

func NewGraphiteClient(address string, interval time.Duration) *GraphiteClient {
	ctx := context.Background()
	client := &GraphiteClient{
		HistogramBucketFormat: "%.2f",
		address:               address,
		metrics:               make(map[string]GraphiteMetric),
		mu:                    new(flu.RWMutex),
		work:                  new(flu.WaitGroup),
	}

	if interval > 0 {
		client.cancel = client.work.Go(ctx, nil, func(ctx context.Context) {
			timer := time.NewTimer(interval)
			defer func() {
				timer.Stop()
				if err := client.Flush(time.Now()); err != nil {
					log.Printf("Failed to flush Graphite metrics: %s", err)
				}
				client.work.Done()
			}()

			for {
				select {
				case <-ctx.Done():
					return
				case now := <-timer.C:
					if err := client.Flush(now); err != nil {
						log.Printf("Failed to flush Graphite metrics: %s", err)
					}
				}
			}
		})
	}

	return client
}

func (g *GraphiteClient) Close() {
	if g.cancel != nil {
		g.cancel()
	}

	g.work.Wait()
}

func (g *GraphiteClient) Flush(now time.Time) error {
	b := new(strings.Builder)
	nowstr := strconv.FormatInt(now.Unix(), 10)

	g.mu.RLock()
	for key, metric := range g.metrics {
		metric.Write(b, nowstr, key)
	}

	g.mu.RUnlock()
	if b.Len() == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), GraphiteTimeout)
	defer cancel()
	data := &flu.PlainText{b.String()}
	conn := flu.Conn{Context: ctx, Network: "tcp", Address: g.address}
	if err := flu.EncodeTo(data, conn); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func (g *GraphiteClient) WithPrefix(prefix string) Registry {
	client := *g
	client.prefix = withPrefix(client.prefix, prefix, ".")
	return &client
}

func (g *GraphiteClient) Counter(name string, labels Labels) Counter {
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

func (g *GraphiteClient) Gauge(name string, labels Labels) Gauge {
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

func (g *GraphiteClient) Histogram(name string, labels Labels, buckets []float64) Histogram {
	key := g.makeKey(name, labels)

	g.mu.RLock()
	entry, ok := g.metrics[key]
	g.mu.RUnlock()

	if !ok {
		g.mu.Lock()
		entry, ok = g.metrics[key]
		if !ok {
			buckets := append(buckets, math.MaxFloat64)
			counters := make([]*GraphiteCounter, len(buckets))
			for i := range buckets {
				counters[i] = new(GraphiteCounter)
			}

			entry = GraphiteHistogram{
				buckets:  buckets,
				counters: counters,
				hbf:      g.HistogramBucketFormat,
			}

			g.metrics[key] = entry
		}
		g.mu.Unlock()
	}

	return entry.(GraphiteHistogram)
}

func (g *GraphiteClient) makeKey(name string, labels Labels) string {
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
