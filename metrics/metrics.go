package metrics

import (
	"strings"

	"github.com/pkg/errors"
)

type Registry interface {
	WithPrefix(prefix string) Registry
	Counter(name string, labels Labels) Counter
	Gauge(name string, labels Labels) Gauge
	Histogram(name string, labels Labels, buckets []float64) Histogram
}

type Counter interface {
	Inc()
	Add(delta float64)
}

type Gauge interface {
	Set(value float64)
	Inc()
	Dec()
	Add(delta float64)
	Sub(delta float64)
}

type Histogram interface {
	Observe(value float64)
}

type Labels []string

func (l Labels) Append(values ...string) Labels {
	return Labels(append(l, values...))
}

func (l Labels) checkLength() {
	if len(l)%2 != 0 {
		panic(errors.Errorf("invalid Labels length: %d", len(l)))
	}
}

func withPrefix(base, suffix, sep string) string {
	if base != "" {
		base += sep
	}

	return base + suffix
}

func (l Labels) Keys() []string {
	if l == nil {
		return nil
	}
	l.checkLength()
	keys := make([]string, len(l)/2)
	for i := 0; i < len(l)/2; i++ {
		keys[i] = l[2*i]
	}
	return keys
}

func (l Labels) Path(sep, esc string) string {
	if l == nil {
		return ""
	}
	l.checkLength()
	b := new(strings.Builder)
	first := true
	for i := 0; i < len(l)/2; i++ {
		if first {
			first = false
		} else {
			b.WriteRune('.')
		}

		b.WriteString(strings.Replace(l[2*i+1], sep, esc, -1))
	}

	return b.String()
}

func (l Labels) Map() map[string]string {
	if l == nil {
		return nil
	}
	l.checkLength()
	m := make(map[string]string, len(l)/2)
	for i := 0; i < len(l)/2; i++ {
		m[l[2*i]] = l[2*i+1]
	}
	return m
}
