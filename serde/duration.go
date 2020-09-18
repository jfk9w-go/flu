package serde

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

type Duration struct {
	time.Duration
}

func (d *Duration) FromString(str string) error {
	var err error
	d.Duration, err = time.ParseDuration(str)
	return err
}

func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
	return d.FromString(node.Value)
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return errors.Wrap(err, "unmarshal string")
	}

	return d.FromString(str)
}
