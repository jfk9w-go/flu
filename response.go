package flu

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Response is a fluent response wrapper.
type Response struct {
	// Error contains an error in case of a request processing error
	// or nil in case of success.
	Error error
	resp  *http.Response
}

// ReadResponseFunc allows to process a http.Response entirely.
type ReadResponseFunc func(*http.Response) error

// ReadResponseFunc executes a ReadResponseFunc if there was no previous error.
func (r *Response) ReadResponseFunc(rf ReadResponseFunc) *Response {
	if r.Error != nil {
		return r
	}

	r.Error = rf(r.resp)
	return r
}

// StatusCodes checks if a http.Response matches a status code from statusCodes.
func (r *Response) StatusCodes(statusCodes ...int) *Response {
	return r.ReadResponseFunc(func(resp *http.Response) error {
		for _, expectedStatusCode := range statusCodes {
			if expectedStatusCode == resp.StatusCode {
				return nil
			}
		}

		return errors.New(resp.Status)
	})
}

// ReadBodyFunc allows to process a http.Response body
type ReadBodyFunc func(io.Reader) error

// ReadBodyFunc executes a ReadBodyFunc.
// It closes the response body after the processing is done.
func (r *Response) ReadBodyFunc(bf ReadBodyFunc) *Response {
	return r.ReadResponseFunc(func(resp *http.Response) error {
		err := bf(resp.Body)
		_ = resp.Body.Close()
		return err
	})
}

// ReadBody checks the content type of the response and reads the body.
func (r *Response) ReadBody(body BodyReader) *Response {
	return r.ReadResponseFunc(func(resp *http.Response) error {
		contentType := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, body.ContentType()) {
			return fmt.Errorf("invalid content type: %s, expected: %s", contentType, body.ContentType())
		}

		err := body.Read(resp.Body)
		_ = resp.Body.Close()

		return err
	})
}

// ReadBytesFunc is a response body byte array processor.
type ReadBytesFunc func([]byte) error

// ReadBytesFunc executes a ReadBytesFunc.
// It closes the response body after the body content has been Read.
func (r *Response) ReadBytesFunc(bf ReadBytesFunc) *Response {
	return r.ReadResponseFunc(func(resp *http.Response) error {
		data, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return err
		}

		return bf(data)
	})
}

// ReadResource allows to save a http.Response body to a WriteResource as is.
func (r *Response) ReadResource(resource WriteResource) *Response {
	return r.ReadBodyFunc(func(body io.Reader) error {
		writer, err := resource.Write()
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, body)
		_ = writer.Close()
		return err
	})
}
