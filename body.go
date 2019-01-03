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
	write(io.Writer) error
	contentType() string
}

// BodyReader is a response body reader.
type BodyReader interface {
	read(io.Reader) error
	contentType() string
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

func (b FormBody) write(body io.Writer) (err error) {
	_, err = io.WriteString(body, url.Values(b).Encode())
	return
}

func (b FormBody) contentType() string {
	return "application/x-www-form-urlencoded"
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

func (b *MultipartFormBody) write(body io.Writer) (err error) {
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

func (b *MultipartFormBody) contentType() string {
	return "multipart/form-data; boundary=" + b.boundary
}

type BufferedBody struct {
	value     interface{}
	marshal   func(interface{}) ([]byte, error)
	unmarshal func([]byte, interface{}) error
	ctype     string
}

func (b *BufferedBody) write(body io.Writer) error {
	data, err := b.marshal(b.value)
	if err != nil {
		return err
	}

	_, err = body.Write(data)
	return err
}

func (b *BufferedBody) contentType() string {
	return b.ctype
}

func (b *BufferedBody) read(body io.Reader) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	return b.unmarshal(data, b.value)
}

// JSON creates an JSON body for a value.
func JSON(value interface{}) *BufferedBody {
	return &BufferedBody{
		value:     value,
		marshal:   json.Marshal,
		unmarshal: json.Unmarshal,
		ctype:     "application/json",
	}
}

// XML creates an XML body for a value.
func XML(value interface{}) *BufferedBody {
	return &BufferedBody{
		value:     value,
		marshal:   xml.Marshal,
		unmarshal: xml.Unmarshal,
		ctype:     "application/xml",
	}
}

func PlainText(text string) *BufferedBody {
	return &BufferedBody{
		value: &text,
		marshal: func(value interface{}) ([]byte, error) {
			return []byte(*value.(*string)), nil
		},
		unmarshal: func(data []byte, value interface{}) error {
			*value.(*string) = string(data)
			return nil
		},
		ctype: "text/plain",
	}
}
