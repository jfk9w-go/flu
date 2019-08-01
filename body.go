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

	"github.com/google/go-querystring/query"
)

type Body interface {
	ContentType() string
}

// BodyWriter is a request body writer.
type BodyWriter interface {
	Body
	Write(io.Writer) error
}

// BodyReader is a response body reader.
type BodyReader interface {
	Body
	Read(io.Reader) error
}

// FormBody represents a form body.
type FormBody struct {
	value     interface{}
	values    url.Values
	multipart *MultipartFormBody
}

// Form creates an empty form.
func Form(value ...interface{}) *FormBody {
	form := new(FormBody)
	if len(value) == 1 {
		form.value = value[0]
	}

	return form
}

// FormValues creates a form with initial values.
func FormValues(values url.Values) *FormBody {
	return &FormBody{
		values: values,
	}
}

func (b *FormBody) Set(key, value string) *FormBody {
	if b.values == nil {
		b.values = url.Values{}
	}

	b.values.Set(key, value)
	return b
}

// Add adds a key-value pair to the form.
func (b *FormBody) Add(key, value string) *FormBody {
	if b.values == nil {
		b.values = url.Values{}
	}

	b.values.Add(key, value)
	return b
}

// AddAll adds a key with multiple values to the form.
func (b *FormBody) AddAll(key string, values ...string) *FormBody {
	for _, value := range values {
		b.Add(key, value)
	}

	return b
}

func (b *FormBody) AddValues(values url.Values) *FormBody {
	for k, vs := range values {
		for i, v := range vs {
			if i == 0 {
				b.values.Set(k, v)
			} else {
				b.values.Add(k, v)
			}
		}
	}

	return b
}

func (b *FormBody) Multipart() *MultipartFormBody {
	if b.multipart == nil {
		b.multipart = MultipartFormFrom(b)
	}

	return b.multipart
}

func (b *FormBody) ContentType() string {
	return "application/x-www-form-urlencoded"
}

func (b *FormBody) Write(body io.Writer) error {
	values, err := b._values()
	if err != nil {
		return err
	}

	_, err = io.WriteString(body, values.Encode())
	return err
}

func (b *FormBody) _values() (url.Values, error) {
	if b.value != nil {
		values, err := query.Values(b.value)
		if err != nil {
			return nil, err
		}

		for k, vs := range b.values {
			for i, v := range vs {
				if i == 0 {
					values.Set(k, v)
				} else {
					values.Add(k, v)
				}
			}
		}

		return values, nil
	} else if b.values != nil {
		return b.values, nil
	}

	return nil, nil
}

// MultipartFormBody represents a form sent as multipart/form-data.
type MultipartFormBody struct {
	*FormBody
	resources map[string]ReadResource
	boundary  string
}

// MultipartForm creates an empty multipart form.
func MultipartForm(values ...interface{}) *MultipartFormBody {
	return MultipartFormFrom(Form(values...))
}

func MultipartFormFrom(formBody *FormBody) *MultipartFormBody {
	return &MultipartFormBody{
		FormBody:  formBody,
		resources: make(map[string]ReadResource),
		boundary:  randomBoundary(),
	}
}

// MultipartFormValues creates a multipart form with initial values.
func MultipartFormValues(values url.Values) *MultipartFormBody {
	return MultipartFormFrom(FormValues(values))
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
	b.FormBody.Add(key, value)
	return b
}

func (b *MultipartFormBody) AddAll(key string, values ...string) *MultipartFormBody {
	b.FormBody.AddAll(key, values...)
	return b
}

func (b *MultipartFormBody) AddValues(values url.Values) *MultipartFormBody {
	b.FormBody.AddValues(values)
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

func (b *MultipartFormBody) Write(body io.Writer) error {
	values, err := b._values()
	if err != nil {
		return err
	}

	writer := multipart.NewWriter(body)
	err = writer.SetBoundary(b.boundary)
	if err != nil {
		return err
	}

	for key, resource := range b.resources {
		var part io.Writer
		part, err = writer.CreateFormFile(key, key)
		if err != nil {
			return err
		}

		var reader io.ReadCloser
		reader, err = resource.Read()
		if err != nil {
			return err
		}

		_, err = io.Copy(part, reader)
		_ = reader.Close()
		if err != nil {
			return err
		}
	}

	for k, vs := range values {
		for _, value := range vs {
			err = writer.WriteField(k, value)
			if err != nil {
				return err
			}
		}
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	return err
}

type (
	// MarshalFunc is a function used to marshal a value to byte array.
	MarshalFunc = func(interface{}) ([]byte, error)

	// UnmarshalFunc is a function used to unmarshal a value from byte array.
	UnmarshalFunc = func([]byte, interface{}) error
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
