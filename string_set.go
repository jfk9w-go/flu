package flu

import (
	"strings"

	"github.com/pkg/errors"
)

type StringSet map[string]bool

func (s StringSet) Has(key string) bool {
	return s[key]
}

func (s StringSet) Add(key string) {
	s[key] = true
}

func (s StringSet) MarshalJSON() ([]byte, error) {
	var b strings.Builder
	b.WriteRune('[')
	first := true
	for value := range s {
		if first {
			first = false
		} else {
			b.WriteString(",")
		}

		b.WriteRune('"')
		b.WriteString(value)
		b.WriteRune('"')
	}

	b.WriteRune(']')
	return []byte(b.String()), nil
}

func (s StringSet) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str[0] != '[' {
		return errors.New("expected array start")
	}

	var b strings.Builder
	var write bool
	for _, c := range str[1:] {
		if c == '"' {
			write = !write
			if !write {
				s.Add(b.String())
				b.Reset()
			}
		} else if write {
			b.WriteRune(c)
		}
	}

	return nil
}

func (s StringSet) Copy() StringSet {
	copy := make(StringSet, len(s))
	for value := range s {
		copy.Add(value)
	}

	return copy
}
