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
	value interface{}
}

func (body jsonBody) EncodeTo(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(body.value)
}

func (body jsonBody) DecodeFrom(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(body.value)
}

func (body jsonBody) ContentType() string {
	return "application/json"
}

func JSON(value interface{}) BodyCodec {
	return jsonBody{value}
}

type xmlBody struct {
	value interface{}
}

func (body xmlBody) EncodeTo(writer io.Writer) error {
	return xml.NewEncoder(writer).Encode(body.value)
}

func (body xmlBody) DecodeFrom(reader io.Reader) error {
	return xml.NewDecoder(reader).Decode(body.value)
}

func (xmlBody) ContentType() string {
	return "application/xml"
}

func XML(value interface{}) BodyCodec {
	return xmlBody{value}
}

type PlainTextBody struct {
	Value string
}

func (body *PlainTextBody) EncodeTo(w io.Writer) error {
	_, err := io.WriteString(w, body.Value)
	return err
}

func (body *PlainTextBody) DecodeFrom(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	body.Value = string(data)
	return nil
}

func (*PlainTextBody) ContentType() string {
	return "text/plain"
}

func PlainText(v string) *PlainTextBody {
	return &PlainTextBody{v}
}
