package metrics

import (
	"strings"

	"github.com/pkg/errors"
)

type Client interface {
	WithPrefix(prefix string) Client
	Counter(name string, labels Labels) Counter
	Gauge(name string, labels Labels) Gauge
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

type Labels []string

func (l Labels) checkLength() {
	if len(l)%2 != 0 {
		panic(errors.Errorf("invalid Labels length: %d", len(l)))
	}
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
	values := make([]string, len(l)/2)
	for i := 0; i < len(l)/2; i++ {
		values[i] = strings.Replace(l[2*i+1], sep, esc, -1)
	}
	return strings.Join(values, ".")
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
