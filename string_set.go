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

func (s StringSet) ForEach(fun func(key string) bool) {
	for k, v := range s {
		if !v {
			continue
		}

		if !fun(k) {
			return
		}
	}
}

func (s StringSet) Delete(key string) {
	s[key] = false
}

func (s StringSet) MarshalJSON() ([]byte, error) {
	var b strings.Builder
	b.WriteRune('[')
	first := true
	s.ForEach(func(key string) bool {
		if first {
			first = false
		} else {
			b.WriteString(", ")
		}

		b.WriteRune('"')
		b.WriteString(key)
		b.WriteRune('"')
		return true
	})

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
				s[b.String()] = true
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
	s.ForEach(func(key string) bool {
		copy.Add(key)
		return true
	})

	return copy
}
