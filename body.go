package flu

import (
	"crypto/rand"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
)

// BodyWriter is a request body writer.
type BodyWriter interface {
	Write(io.Writer) error
	ContentType() string
}

// BodyReader is a response body reader.
type BodyReader interface {
	Read(io.Reader) error
	ContentType() string
}

// FormBody represents a form body.
type FormBody url.Values

// Form creates an empty form.
func Form() FormBody {
	return FormWith(url.Values{})
}

// FormWith creates a form with initial values.
func FormWith(values url.Values) FormBody {
	return FormBody(values)
}

// Add adds a key-value pair to the form.
func (b FormBody) Add(key, value string) FormBody {
	url.Values(b).Add(key, value)
	return b
}

// AddAll adds a key with multiple values to the form.
func (b FormBody) AddAll(key string, values ...string) FormBody {
	for _, value := range values {
		b.Add(key, value)
	}

	return b
}

func (b FormBody) ContentType() string {
	return "application/x-www-form-urlencoded"
}

func (b FormBody) Write(body io.Writer) (err error) {
	_, err = io.WriteString(body, url.Values(b).Encode())
	return
}

// MultipartFormBody represents a form sent as multipart/form-data.
type MultipartFormBody struct {
	values    url.Values
	resources map[string]ReadResource
	boundary  string
}

// MultipartForm creates an empty multipart form.
func MultipartForm() *MultipartFormBody {
	return MultipartFormWith(url.Values{})
}

// MultipartFormWith creates a multipart form with initial values.
func MultipartFormWith(values url.Values) *MultipartFormBody {
	return &MultipartFormBody{
		values:    values,
		resources: make(map[string]ReadResource),
		boundary:  randomBoundary(),
	}
}

func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

// Add adds a key-value pair to the form.
func (b *MultipartFormBody) Add(key, value string) *MultipartFormBody {
	b.values.Add(key, value)
	return b
}

// Resource adds a Resource to the form.
func (b *MultipartFormBody) Resource(key string, resource ReadResource) *MultipartFormBody {
	b.resources[key] = resource
	return b
}

func (b *MultipartFormBody) ContentType() string {
	return "multipart/form-data; boundary=" + b.boundary
}

func (b *MultipartFormBody) Write(body io.Writer) (err error) {
	var writer = multipart.NewWriter(body)
	err = writer.SetBoundary(b.boundary)
	if err != nil {
		return
	}

	for key, resource := range b.resources {
		var part io.Writer
		part, err = writer.CreateFormFile(key, key)
		if err != nil {
			return
		}

		var reader io.ReadCloser
		reader, err = resource.Read()
		if err != nil {
			return
		}

		_, err = io.Copy(part, reader)
		_ = reader.Close()
		if err != nil {
			return
		}
	}

	for key, values := range b.values {
		for _, value := range values {
			err = writer.WriteField(key, value)
			if err != nil {
				return
			}
		}
	}

	err = writer.Close()
	if err != nil {
		return
	}

	return
}

type (
	// MarshalFunc is a function used to marshal a value to byte array.
	MarshalFunc func(interface{}) ([]byte, error)

	// UnmarshalFunc is a function used to unmarshal a value from byte array.
	UnmarshalFunc func([]byte, interface{}) error
)

// BodyReadWriter is a generic request body writer and response body reader.
// A value is marshalled to and unmarshalled from a byte array.
type BodyReadWriter struct {
	contentType   string
	marshalFunc   MarshalFunc
	unmarshalFunc UnmarshalFunc
	value         interface{}
}

// ReadWriteBody creates a new instance of a BodyReadWriter.
func ReadWriteBody(contentType string, marshalFunc MarshalFunc, unmarshalFunc UnmarshalFunc, value interface{}) *BodyReadWriter {
	return &BodyReadWriter{contentType, marshalFunc, unmarshalFunc, value}
}

func (b *BodyReadWriter) ContentType() string {
	return b.contentType
}

func (b *BodyReadWriter) Write(body io.Writer) error {
	data, err := b.marshalFunc(b.value)
	if err != nil {
		return err
	}

	_, err = body.Write(data)
	return err
}

func (b *BodyReadWriter) Read(body io.Reader) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	return b.unmarshalFunc(data, b.value)
}

// JSON creates an JSON body for a value.
func JSON(value interface{}) *BodyReadWriter {
	return &BodyReadWriter{
		contentType:   "application/json",
		value:         value,
		marshalFunc:   json.Marshal,
		unmarshalFunc: json.Unmarshal,
	}
}

// XML creates an XML body for a value.
func XML(value interface{}) *BodyReadWriter {
	return &BodyReadWriter{
		contentType:   "application/xml",
		value:         value,
		marshalFunc:   xml.Marshal,
		unmarshalFunc: xml.Unmarshal,
	}
}

// PlainText creates a plain/text body.
func PlainText(value *string) *BodyReadWriter {
	return &BodyReadWriter{
		contentType: "text/plain",
		marshalFunc: func(interface{}) ([]byte, error) {
			return []byte(*value), nil
		},
		unmarshalFunc: func(data []byte, _ interface{}) error {
			*value = string(data)
			return nil
		},
	}
}
