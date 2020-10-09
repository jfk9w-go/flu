package serde

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

var TimeLayout = "2006-01-02 15:04:05"

type Time struct {
	time.Time
}

func (t Time) String() string {
	return t.Time.Format(TimeLayout)
}

func (t *Time) FromString(str string) error {
	value, err := time.Parse(TimeLayout, str)
	if err != nil {
		return errors.Wrap(err, "parse time")
	}
	t.Time = value
	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	return t.FromString(str)
}

func (t *Time) MarshalYAML() (interface{}, error) {
	return t.String(), nil
}

func (t *Time) UnmarshalYAML(node *yaml.Node) error {
	return t.FromString(node.Value)
}
