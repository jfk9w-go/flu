package csv

import (
	"time"

	"github.com/pkg/errors"
)

type RowOutput interface {
	Output(row *Row) error
}

type Row struct {
	Values map[string]string
	Err    error
}

func (r *Row) Get(key string, value Value) {
	if r.Err != nil {
		return
	}

	rawValue, ok := r.Values[key]
	if !ok {
		r.Err = errors.Errorf("index '%s' out of bounds", key)
		return
	}

	if err := value.Parse(rawValue); err != nil {
		r.Err = errors.Wrapf(err, "at column '%s' (row: '%+v')", key, r.Values)
	}
}

func (r *Row) Time(key string, layout string, location *time.Location) time.Time {
	v := Time{
		Layout:   layout,
		Location: location,
	}

	r.Get(key, &v)
	return v.Value
}

func (r *Row) Float(key string, bitSize int) float64 {
	v := Float{BitSize: bitSize}
	r.Get(key, &v)
	return v.Value
}

func (r *Row) String(key string) string {
	var v String
	r.Get(key, &v)
	return v.Value()
}

func (r *Row) Uint(key string, base, bitSize int) uint64 {
	v := Uint{
		Base:    base,
		BitSize: bitSize,
	}

	r.Get(key, &v)
	return v.Value
}
