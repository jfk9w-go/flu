package flu

import (
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
	Error    error
	httpResp *http.Response
}

// ReadResponseFunc allows to process a http.Response entirely.
type ReadResponseFunc func(*http.Response) error

// ReadResponseFunc executes a ReadResponseFunc if there was no previous error.
func (resp *Response) ReadResponseFunc(readResponse ReadResponseFunc) *Response {
	if resp.Error != nil {
		return resp
	}

	resp.Error = readResponse(resp.httpResp)
	return resp
}

// CheckStatusCode checks if a http.Response matches a status code from statusCodes.
func (resp *Response) CheckStatusCode(allowedStatusCodes ...int) *Response {
	return resp.ReadResponseFunc(func(httpResp *http.Response) error {
		for _, expectedStatusCode := range allowedStatusCodes {
			if expectedStatusCode == httpResp.StatusCode {
				return nil
			}
		}

		return fmt.Errorf("invalid status code: %d", httpResp.StatusCode)
	})
}

// ReadBodyFunc allows to process a http.Response body
type ReadBodyFunc func(io.Reader) error

// ReadBodyFunc executes a ReadBodyFunc.
// It closes the response body after the processing is done.
func (resp *Response) ReadBodyFunc(readBody ReadBodyFunc) *Response {
	return resp.ReadResponseFunc(func(resp *http.Response) error {
		err := readBody(resp.Body)
		_ = resp.Body.Close()
		return err
	})
}

// ReadBody checks the content type of the response and reads the body.
func (resp *Response) ReadBody(reader BodyReader) *Response {
	return resp.ReadResponseFunc(func(resp *http.Response) error {
		contentType := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, reader.ContentType()) {
			return fmt.Errorf("invalid content type: %s", contentType)
		}

		err := reader.Read(resp.Body)
		_ = resp.Body.Close()

		return err
	})
}

// ReadBytesFunc is a response body byte array processor.
type ReadBytesFunc func([]byte) error

// ReadBytesFunc executes a ReadBytesFunc.
// It closes the response body after the body content has been Read.
func (resp *Response) ReadBytesFunc(readBytes ReadBytesFunc) *Response {
	return resp.ReadResponseFunc(func(httpResp *http.Response) error {
		data, err := ioutil.ReadAll(httpResp.Body)
		_ = httpResp.Body.Close()
		if err != nil {
			return err
		}

		return readBytes(data)
	})
}

// ReadResource allows to save a http.Response body to a WriteResource as is.
func (resp *Response) ReadResource(resource WriteResource) *Response {
	return resp.ReadBodyFunc(func(body io.Reader) error {
		writer, err := resource.Write()
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, body)
		_ = writer.Close()
		return err
	})
}
