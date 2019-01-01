package flu

import (
	"crypto/rand"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

// RequestBodyBuilder is a request body.
type RequestBodyBuilder interface {
	build(io.Writer) error
	contentType() string
}

// FormBodyBuilder represents a form body.
type FormBodyBuilder url.Values

// Form creates an empty form.
func Form() FormBodyBuilder {
	return FormWith(url.Values{})
}

// FormWith creates a form with initial values.
func FormWith(values url.Values) FormBodyBuilder {
	return FormBodyBuilder(values)
}

// Add adds a key-value pair to the form.
func (b FormBodyBuilder) Add(key, value string) FormBodyBuilder {
	url.Values(b).Add(key, value)
	return b
}

// AddAll adds a key with multiple values to the form.
func (b FormBodyBuilder) AddAll(key string, values ...string) FormBodyBuilder {
	for _, value := range values {
		b.Add(key, value)
	}

	return b
}

func (b FormBodyBuilder) build(body io.Writer) (err error) {
	_, err = io.WriteString(body, url.Values(b).Encode())
	return
}

func (b FormBodyBuilder) contentType() string {
	return "application/x-www-form-urlencoded"
}

// MultipartFormBodyBuilder represents a form sent as multipart/form-data.
type MultipartFormBodyBuilder struct {
	values    url.Values
	resources map[string]ReadResource
	boundary  string
}

// MultipartForm creates an empty multipart form.
func MultipartForm() MultipartFormBodyBuilder {
	return MultipartFormWith(url.Values{})
}

// MultipartFormWith creates a multipart form with initial values.
func MultipartFormWith(values url.Values) MultipartFormBodyBuilder {
	return MultipartFormBodyBuilder{
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
func (b MultipartFormBodyBuilder) Add(key, value string) MultipartFormBodyBuilder {
	b.values.Add(key, value)
	return b
}

// Resource adds a Resource to the form.
func (b MultipartFormBodyBuilder) Resource(key string, resource ReadResource) MultipartFormBodyBuilder {
	b.resources[key] = resource
	return b
}

func (b MultipartFormBodyBuilder) build(body io.Writer) (err error) {
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

func (b MultipartFormBodyBuilder) contentType() string {
	return "multipart/form-data; boundary=" + b.boundary
}

// JsonBodyBuilder represents a request body sent as a JSON.
type JsonBodyBuilder struct {
	value interface{}
}

// Json creates a request JSON body from a value.
func Json(value interface{}) *JsonBodyBuilder {
	return &JsonBodyBuilder{value}
}

func (b JsonBodyBuilder) build(body io.Writer) error {
	var data, err = json.Marshal(b.value)
	if err != nil {
		return err
	}

	_, err = body.Write(data)
	return err
}

func (b JsonBodyBuilder) contentType() string {
	return "application/json"
}

type XmlBodyBuilder struct {
	value interface{}
}

// Xml creates a request XML body from a value.
func Xml(value interface{}) XmlBodyBuilder {
	return XmlBodyBuilder{value}
}

func (b XmlBodyBuilder) build(body io.Writer) error {
	var data, err = xml.Marshal(b.value)
	if err != nil {
		return err
	}

	_, err = body.Write(data)
	return err
}

func (b XmlBodyBuilder) contentType() string {
	return "application/xml"
}
