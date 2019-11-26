package flu

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
)

type EncoderTo interface {
	EncodeTo(io.Writer) error
}

type DecoderFrom interface {
	DecodeFrom(io.Reader) error
}

type Body interface {
	ContentType() string
}

type BodyEncoderTo interface {
	Body
	EncoderTo
}

type BodyDecoderFrom interface {
	Body
	DecoderFrom
}

type BodyCodec interface {
	Body
	EncoderTo
	DecoderFrom
}

type jsonBody struct {
	v interface{}
}

func (b jsonBody) EncodeTo(w io.Writer) error {
	return json.NewEncoder(w).Encode(b.v)
}

func (b jsonBody) DecodeFrom(r io.Reader) error {
	return json.NewDecoder(r).Decode(b.v)
}

func (b jsonBody) ContentType() string {
	return "application/json"
}

func JSON(v interface{}) BodyCodec {
	return jsonBody{v}
}

type xmlBody struct {
	v interface{}
}

func (b xmlBody) EncodeTo(w io.Writer) error {
	return xml.NewEncoder(w).Encode(b.v)
}

func (b xmlBody) DecodeFrom(r io.Reader) error {
	return xml.NewDecoder(r).Decode(b.v)
}

func (b xmlBody) ContentType() string {
	return "application/xml"
}

func XML(v interface{}) BodyCodec {
	return xmlBody{v}
}

type PlainTextBody struct {
	Value string
}

func (c PlainTextBody) EncodeTo(w io.Writer) error {
	_, err := io.WriteString(w, c.Value)
	return err
}

func (c PlainTextBody) DecodeFrom(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	p := &c.Value
	*p = string(data)
	return nil
}

func (c PlainTextBody) ContentType() string {
	return "text/plain"
}

func PlainText(v string) PlainTextBody {
	return PlainTextBody{v}
}
