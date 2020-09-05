package metrics

import "log"

type DummyCounter struct {
	Name   string
	Labels Labels
	Log    bool
}

func (c DummyCounter) Inc() {
	if c.Log {
		log.Printf("metrics.DummyCounter %s %v inc", c.Name, c.Labels)
	}
}

func (c DummyCounter) Add(delta float64) {
	if c.Log {
		log.Printf("metrics.DummyCounter %s %v add %.4f", c.Name, c.Labels, delta)
	}
}

type DummyGauge struct {
	Name   string
	Labels Labels
	Log    bool
}

func (g DummyGauge) Set(value float64) {
	if g.Log {
		log.Printf("metrics.DummyGauge %s %v set %.4f", g.Name, g.Labels, value)
	}
}

func (g DummyGauge) Inc() {
	if g.Log {
		log.Printf("metrics.DummyGauge %s %v inc", g.Name, g.Labels)
	}
}

func (g DummyGauge) Dec() {
	if g.Log {
		log.Printf("metrics.DummyGauge %s %v dec", g.Name, g.Labels)
	}
}

func (g DummyGauge) Add(delta float64) {
	if g.Log {
		log.Printf("metrics.DummyGauge %s %v add %.4f", g.Name, g.Labels, delta)
	}
}

func (g DummyGauge) Sub(delta float64) {
	if g.Log {
		log.Printf("metrics.DummyGauge %s %v sub %.4f", g.Name, g.Labels, delta)
	}
}

type DummyRegistry struct {
	Prefix string
	Log    bool
}

func (d DummyRegistry) WithPrefix(prefix string) Registry {
	d.Prefix = withPrefix(d.Prefix, prefix, ".")
	return d
}

func (d DummyRegistry) Counter(name string, labels Labels) Counter {
	if d.Prefix != "" {
		d.Prefix += "."
	}

	return DummyCounter{
		Name:   d.Prefix + name,
		Labels: labels,
		Log:    d.Log,
	}
}

func (d DummyRegistry) Gauge(name string, labels Labels) Gauge {
	if d.Prefix != "" {
		d.Prefix += "."
	}

	return DummyGauge{
		Name:   d.Prefix + name,
		Labels: labels,
		Log:    d.Log,
	}
}
