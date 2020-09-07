package flu

import (
	"encoding/json"
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

func (s Uint64Set) Delete(key uint64) {
	delete(s, key)
}

type BigUint64Set struct {
	Avg string   `json:"a,omitempty"`
	Dev []string `json:"d,omitempty"`
}

func (s Uint64Set) MarshalJSON() ([]byte, error) {
	var sum uint64
	for value := range s {
		sum += value
	}

	size := len(s)
	repr := BigUint64Set{Dev: make([]string, size)}
	var avg uint64
	if size > 0 {
		avg = sum / uint64(size)
		repr.Avg = strconv.FormatUint(avg, 36)
	}

	i := 0
	for value := range s {
		repr.Dev[i] = strconv.FormatUint(value-avg, 36)
	}

	return json.Marshal(repr)
}

func (s Uint64Set) UnmarshalJSON(data []byte) error {
	var repr BigUint64Set
	if err := json.Unmarshal(data, &repr); err != nil {
		return errors.Wrap(err, "unmarshal repr")
	}

	if len(repr.Dev) == 0 {
		return nil
	}

	avg, err := strconv.ParseUint(repr.Avg, 36, 64)
	if err != nil {
		return errors.Wrap(err, "parse avg")
	}

	for _, devstr := range repr.Dev {
		dev, err := strconv.ParseUint(devstr, 36, 64)
		if err != nil {
			return errors.Wrapf(err, "parse dev: %s", devstr)
		}

		s[avg+dev] = true
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
