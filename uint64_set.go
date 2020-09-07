package flu

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/pkg/errors"
)

type Uint64Set map[uint64]bool

func (s Uint64Set) Has(key uint64) bool {
	return s[key]
}

func (s Uint64Set) Add(key uint64) {
	s[key] = true
}

type BigUint64Set struct {
	Base string    `json:"b,omitempty"`
	Off  StringSet `json:"o,omitempty"`
}

func (s Uint64Set) MarshalJSON() ([]byte, error) {
	var base uint64 = math.MaxUint64
	for value := range s {
		if value < base {
			base = value
		}
	}

	size := len(s)
	repr := BigUint64Set{Off: make(StringSet, size)}
	if size > 0 {
		repr.Base = strconv.FormatUint(base, 36)
	}

	for value := range s {
		off := value - base
		if off == 0 {
			continue
		}

		repr.Off.Add(strconv.FormatUint(off, 36))
	}

	return json.Marshal(repr)
}

func (s Uint64Set) UnmarshalJSON(data []byte) error {
	repr := BigUint64Set{Off: make(StringSet)}
	if err := json.Unmarshal(data, &repr); err != nil {
		return errors.Wrap(err, "unmarshal repr")
	}

	if string(data) == "{}" {
		return nil
	}

	base, err := strconv.ParseUint(repr.Base, 36, 64)
	if err != nil {
		return errors.Wrap(err, "parse avg")
	}

	s.Add(base)
	for str := range repr.Off {
		off, err := strconv.ParseUint(str, 36, 64)
		if err != nil {
			return errors.Wrapf(err, "parse offset: %s", str)
		}

		s.Add(base + off)
	}

	return nil
}

func (s Uint64Set) Copy() Uint64Set {
	copy := make(Uint64Set, len(s))
	for value := range s {
		copy.Add(value)
	}

	return copy
}
