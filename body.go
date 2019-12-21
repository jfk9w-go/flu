package flu

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
)

type jsonBody struct {
	value interface{}
}

func (b jsonBody) WriteTo(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(b.value)
}

func (b jsonBody) ReadFrom(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(b.value)
}

func (b jsonBody) ContentType() string {
	return "application/json"
}

func JSON(value interface{}) BodyReadWriter {
	return jsonBody{value}
}

type xmlBody struct {
	value interface{}
}

func (b xmlBody) WriteTo(writer io.Writer) error {
	return xml.NewEncoder(writer).Encode(b.value)
}

func (b xmlBody) ReadFrom(reader io.Reader) error {
	return xml.NewDecoder(reader).Decode(b.value)
}

func (b xmlBody) ContentType() string {
	return "application/xml"
}

func XML(value interface{}) BodyReadWriter {
	return xmlBody{value}
}

type PlainTextBody struct {
	Value string
}

func (b *PlainTextBody) WriteTo(w io.Writer) error {
	_, err := io.WriteString(w, b.Value)
	return err
}

func (b *PlainTextBody) ReadFrom(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	b.Value = string(data)
	return nil
}

func (b *PlainTextBody) ContentType() string {
	return "text/plain"
}

func PlainText(v string) *PlainTextBody {
	return &PlainTextBody{v}
}
