package flu

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

// Response is a fluent response wrapper.
type Response struct {
	// Error contains an error in case of a request processing error
	// or nil in case of success.
	Error error
	resp  *http.Response
}

// ResponseProcessor allows to process a http.Response entirely.
type ResponseProcessor func(*http.Response) error

// ProcessResponse executes a ResponseProcessor if there was no previous error.
func (r *Response) ProcessResponse(processor ResponseProcessor) *Response {
	if r.Error != nil {
		return r
	}

	r.Error = processor(r.resp)
	return r
}

// StatusCodes checks if a http.Response matches a status code from statusCodes.
func (r *Response) StatusCodes(statusCodes ...int) *Response {
	return r.ProcessResponse(func(resp *http.Response) error {
		for _, expectedStatusCode := range statusCodes {
			if expectedStatusCode == resp.StatusCode {
				return nil
			}
		}

		return errors.New(resp.Status)
	})
}

// BodyProcessor allows to process a http.Response body
type BodyProcessor func(io.Reader) error

// ProcessBody executes a BodyProcessor.
// It closes the response body after the processing is done.
func (r *Response) ProcessBody(processor BodyProcessor) *Response {
	return r.ProcessResponse(func(resp *http.Response) error {
		err := processor(resp.Body)
		_ = resp.Body.Close()
		return err
	})
}

// BufferedBodyProcessor allows to process a http.Response body contents as buffered bytes.
type BufferedBodyProcessor func([]byte) error

// ProcessBufferedBody reads the response body to a byte array and executes a BufferedBodyProcessor.
func (r *Response) ProcessBufferedBody(processor BufferedBodyProcessor) *Response {
	return r.ProcessResponse(func(resp *http.Response) error {
		data, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return err
		}

		return processor(data)
	})
}

// ReadJSON allows to parse a http.Response body as JSON.
func (r *Response) ReadJSON(value interface{}) *Response {
	return r.ProcessBufferedBody(func(data []byte) error {
		return json.Unmarshal(data, value)
	})
}

// ReadString allows to parse a http.Response body as a string.
func (r *Response) ReadString(value *string) *Response {
	return r.ProcessBufferedBody(func(data []byte) error {
		*value = string(data)
		return nil
	})
}

// ReadResource allows to save a http.Response body to a WriteResource as is.
func (r *Response) ReadResource(resource WriteResource) *Response {
	return r.ProcessBody(func(body io.Reader) error {
		writer, err := resource.Write()
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, body)
		_ = writer.Close()
		return err
	})
}
