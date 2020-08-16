package flu

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

type JSON struct {
	Value interface{}
}

func (j JSON) EncodeTo(w io.Writer) error {
	return json.NewEncoder(w).Encode(j.Value)
}

func (j JSON) DecodeFrom(r io.Reader) error {
	return json.NewDecoder(r).Decode(j.Value)
}

func (j JSON) ContentType() string {
	return "application/json"
}

type XML struct {
	Value interface{}
}

func (x XML) EncodeTo(w io.Writer) error {
	return xml.NewEncoder(w).Encode(x.Value)
}

func (x XML) DecodeFrom(r io.Reader) error {
	return xml.NewDecoder(r).Decode(x.Value)
}

func (x XML) ContentType() string {
	return "application/xml"
}

type PlainText struct {
	Value string
}

func (t *PlainText) EncodeTo(w io.Writer) error {
	_, err := io.WriteString(w, t.Value)
	return err
}

func (t *PlainText) DecodeFrom(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	t.Value = string(data)
	return nil
}

func (t *PlainText) ContentType() string {
	return "text/plain; charset=utf-8"
}

type YAML struct {
	Value interface{}
}

func (y YAML) EncodeTo(w io.Writer) error {
	enc := yaml.NewEncoder(w)
	defer enc.Close()
	return enc.Encode(y.Value)
}

func (y YAML) DecodeFrom(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(y.Value)
}
