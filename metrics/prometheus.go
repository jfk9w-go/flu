package metrics

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/pkg/errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusKey struct {
	Namespace string
	Name      string
}

type PrometheusListener struct {
	server   *http.Server
	registry *prometheus.Registry
	prefix   string
	entries  map[prometheusKey]interface{}
	mu       *sync.RWMutex
}

func NewPrometheusListener(address string) PrometheusListener {
	u, err := url.Parse(address)
	if err != nil {
		panic(errors.Wrap(err, "invalid prometheus address"))
	}

	mux := http.NewServeMux()
	registry := prometheus.NewRegistry()
	mux.Handle(u.Path, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	server := &http.Server{
		Addr:    u.Host,
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Prometheus listener: %s", err)
		}
	}()

	return PrometheusListener{
		server:   server,
		registry: registry,
		entries:  make(map[prometheusKey]interface{}),
		mu:       new(sync.RWMutex),
	}
}

func (p PrometheusListener) MustRegister(cs ...prometheus.Collector) PrometheusListener {
	p.registry.MustRegister(cs...)
	return p
}

func (p PrometheusListener) Close(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

func (p PrometheusListener) WithPrefix(prefix string) Registry {
	if p.prefix != "" {
		p.prefix += "_" + prefix
	} else {
		p.prefix = prefix
	}

	return p
}

func (p PrometheusListener) Counter(name string, labels Labels) Counter {
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
				p.registry.MustRegister(counter)
				p.entries[key] = counter
				entry = counter
			} else {
				vec := prometheus.NewCounterVec(opts, labels.Keys())
				p.registry.MustRegister(vec)
				p.entries[key] = vec
				entry = vec
			}
		}

		p.mu.Unlock()
	}

	if labels != nil {
		entry = entry.(*prometheus.CounterVec).With(labels.Map())
	}

	return entry.(Counter)
}

func (p PrometheusListener) Gauge(name string, labels Labels) Gauge {
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
				Name:      name,
			}

			if labels == nil {
				gauge := prometheus.NewGauge(opts)
				p.registry.MustRegister(gauge)
				p.entries[key] = gauge
				entry = gauge
			} else {
				vec := prometheus.NewGaugeVec(opts, labels.Keys())
				p.registry.MustRegister(vec)
				p.entries[key] = vec
				entry = vec
			}
		}

		p.mu.Unlock()
	}

	if labels != nil {
		entry = entry.(*prometheus.GaugeVec).With(labels.Map())
	}

	return entry.(Gauge)
}
