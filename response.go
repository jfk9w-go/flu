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

func NewStatusCodeError(r *http.Response) StatusCodeError {
	e := StatusCodeError{Code: r.StatusCode}
	if r.Body != nil {
		defer r.Body.Close()
		text := PlainText("")
		err := text.DecodeFrom(r.Body)
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
	return r.complete(NewStatusCodeError(r.http))
}

// Decode reads the response body.
func (r *Response) Decode(decoder DecoderFrom) *Response {
	if r.Error != nil {
		return r
	}
	return r.complete(DecodeFrom(Xable{R: r.http.Body}, decoder))
}

type ContentTypeError string

func (e ContentTypeError) Error() string {
	return fmt.Sprintf("invalid content type: %s", string(e))
}

// DecodeBody checks the response Content-Type header.
// If there is no match, sets the error to ContentTypeError.
// Otherwise proceeds with reading.
func (r *Response) DecodeBody(decoder BodyDecoderFrom) *Response {
	if r.Error != nil {
		return r
	}
	contentType := r.http.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, decoder.ContentType()) {
		r.Error = ContentTypeError(contentType)
		return r
	}
	return r.Decode(decoder)
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
