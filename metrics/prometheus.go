package metrics

import (
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusKey struct {
	Namespace string
	Name      string
}

type PrometheusClient struct {
	prefix  string
	entries map[prometheusKey]interface{}
	mu      *sync.RWMutex
}

func NewPrometheusClient(address string) PrometheusClient {
	http.Handle("/metrics", promhttp.Handler())
	go func() { log.Fatal(http.ListenAndServe(address, nil)) }()
	return PrometheusClient{
		entries: make(map[prometheusKey]interface{}),
		mu:      new(sync.RWMutex),
	}
}

func (p PrometheusClient) WithPrefix(prefix string) Client {
	if p.prefix != "" {
		p.prefix += "_" + prefix
	} else {
		p.prefix = prefix
	}
	return p
}

func (p PrometheusClient) Counter(name string, labels Labels) Counter {
	key := prometheusKey{p.prefix, name}
	p.mu.RLock()
	entry, ok := p.entries[key]
	p.mu.RUnlock()
	if !ok {
		p.mu.Lock()
		entry, ok = p.entries[key]
		if !ok {
			opts := prometheus.CounterOpts{
				Namespace: p.prefix,
				Name:      name,
			}
			if labels == nil {
				counter := prometheus.NewCounter(opts)
				prometheus.MustRegister(counter)
				p.entries[key] = counter
				entry = counter
			} else {
				vec := prometheus.NewCounterVec(opts, labels.Keys())
				prometheus.MustRegister(vec)
				p.entries[key] = vec
				entry = vec
			}
		}
		p.mu.Unlock()
	}

	if labels != nil {
		entry = entry.(*prometheus.CounterVec).With(prometheus.Labels(labels))
	}

	return entry.(Counter)
}

func (p PrometheusClient) Gauge(name string, labels Labels) Gauge {
	key := prometheusKey{p.prefix, name}
	p.mu.RLock()
	entry, ok := p.entries[key]
	p.mu.RUnlock()
	if !ok {
		p.mu.Lock()
		entry, ok = p.entries[key]
		if !ok {
			opts := prometheus.GaugeOpts{
				Namespace: p.prefix,
				Subsystem: p.prefix,
				Name:      name,
			}
			if labels == nil {
				gauge := prometheus.NewGauge(opts)
				prometheus.MustRegister(gauge)
				p.entries[key] = gauge
				entry = gauge
			} else {
				vec := prometheus.NewGaugeVec(opts, labels.Keys())
				prometheus.MustRegister(vec)
				p.entries[key] = vec
				entry = vec
			}
		}
		p.mu.Unlock()
	}

	if labels != nil {
		entry = entry.(*prometheus.GaugeVec).With(prometheus.Labels(labels))
	}

	return entry.(Gauge)
}

func (p PrometheusClient) Close() error {
	return nil
}
