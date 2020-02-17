package metrics

import (
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GraphiteMetric interface {
	Reset() (float64, bool)
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
	return value, value == zero
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
	work    *sync.WaitGroup
	halt    chan struct{}
	mu      *sync.RWMutex
}

func NewGraphiteClient(address string, interval time.Duration) GraphiteClient {
	g := GraphiteClient{
		address: address,
		metrics: make(map[string]GraphiteMetric),
		work:    new(sync.WaitGroup),
		halt:    make(chan struct{}),
		mu:      new(sync.RWMutex),
	}

	g.work.Add(1)
	go func() {
		defer g.work.Done()
		timer := time.NewTimer(interval)
		for {
			select {
			case <-g.halt:
				return
			case now := <-timer.C:
				g.FlushValues(now)
			}
		}
	}()

	return g
}

func (g GraphiteClient) FlushValues(now time.Time) {
	b := new(strings.Builder)
	nowstr := strconv.FormatInt(now.UnixNano()/1e6, 10)

	g.mu.RLock()
	for key, metric := range g.metrics {
		value, reset := metric.Reset()
		if !reset {
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
		return
	}

	conn, err := net.Dial("tcp", g.address)
	if err != nil {
		log.Printf("Failed to connect to graphite on %s: %s", g.address, err)
		return
	}

	_, err = conn.Write([]byte(b.String()))
	_ = conn.Close()
	if err != nil {
		log.Printf("Failed to write data to graphite on %s: %s", g.address, err)
	}
}

func (g GraphiteClient) Close() error {
	g.halt <- struct{}{}
	g.FlushValues(time.Now())
	g.work.Wait()
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
			entry := new(GraphiteCounter)
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
			entry := new(GraphiteGauge)
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

	values := labels.Values(".", "_")
	prefix += values
	if values != "" {
		prefix += "."
	}

	return prefix + name
}
