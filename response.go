package flu

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Response is a fluent response wrapper.
type Response struct {
	http *http.Response

	// Error contains an error in case of a request processing error
	// or nil in case of success.
	Error error
}

type ResponseHandler interface {
	Handle(*http.Response) error
}

// HandleResponse executes a ResponseHandler if no previous handling errors occurred.
func (r *Response) HandleResponse(handler ResponseHandler) *Response {
	if r.Error != nil {
		return r
	}
	return r.complete(handler.Handle(r.http))
}

type StatusCodeError struct {
	Code int
	Text string
}

func (e StatusCodeError) Error() string {
	text := fmt.Sprintf("%d", e.Code)
	if e.Text != "" {
		text += " (" + e.Text + ")"
	}
	return text
}

func createStatusCodeError(r *http.Response) StatusCodeError {
	e := StatusCodeError{Code: r.StatusCode}
	if r.Body != nil {
		defer r.Body.Close()
		text := PlainText("")
		err := text.ReadFrom(r.Body)
		if err != nil {
			text.Value = fmt.Sprintf("response body read error: %s", err.Error())
		}
		e.Text = text.Value
	}
	return e
}

// CheckStatusCode checks the response status code and sets the error to StatusCodeError if there is no match.
func (r *Response) CheckStatusCode(codes ...int) *Response {
	if r.Error != nil {
		return r
	}
	for _, c := range codes {
		if c == r.http.StatusCode {
			return r
		}
	}
	return r.complete(createStatusCodeError(r.http))
}

// Read reads the response body.
func (r *Response) Read(reader ReaderFrom) *Response {
	if r.Error != nil {
		return r
	}
	return r.complete(Read(Xable{R: r.http.Body}, reader))
}

type ContentTypeError string

func (e ContentTypeError) Error() string {
	return fmt.Sprintf("invalid content type: %s", string(e))
}

// ReadBody checks the response Content-Type header.
// If there is no match, sets the error to ContentTypeError.
// Otherwise proceeds with reading.
func (r *Response) ReadBody(reader BodyReader) *Response {
	if r.Error != nil {
		return r
	}
	contentType := r.http.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, reader.ContentType()) {
		r.Error = ContentTypeError(contentType)
		return r
	}
	return r.Read(reader)
}

type bodyReaderHandler struct {
	reader io.Reader
}

func (b *bodyReaderHandler) Handle(resp *http.Response) error {
	b.reader = resp.Body
	return nil
}

func (r *Response) Reader() (io.Reader, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	h := new(bodyReaderHandler)
	err := r.HandleResponse(h).Error
	return h.reader, err
}

type WritableError struct {
	Err error
}

func (e WritableError) Error() string {
	return fmt.Sprintf("failed to create writer: %s", e.Err)
}

//noinspection GoUnhandledErrorResult
func (r *Response) ReadBodyTo(writable Writable) *Response {
	if r.Error != nil {
		return r
	}
	return r.complete(Copy(Xable{R: r.http.Body}, writable))
}

func (r *Response) complete(err error) *Response {
	r.Error = err
	return r
}
