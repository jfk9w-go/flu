package flu

import (
	"fmt"
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
func (r *Response) HandleResponse(h ResponseHandler) *Response {
	if r.Error != nil {
		return r
	}

	return r.complete(h.Handle(r.http))
}

type StatusCodeError int

func (e StatusCodeError) Error() string {
	return fmt.Sprintf("invalid status code: %d", int(e))
}

// CheckStatusCode checks the response status code and sets the error to StatusCodeError if there is no match.
func (r *Response) CheckStatusCode(cs ...int) *Response {
	if r.Error != nil {
		return r
	}

	for _, c := range cs {
		if c == r.http.StatusCode {
			return r
		}
	}

	return r.complete(StatusCodeError(r.http.StatusCode))
}

// Decode decodes the response body.
func (r *Response) Decode(d DecoderFrom) *Response {
	if r.Error != nil {
		return r
	}

	body := r.http.Body
	defer body.Close()
	return r.complete(d.DecodeFrom(body))
}

type ContentTypeError string

func (e ContentTypeError) Error() string {
	return fmt.Sprintf("invalid content type: %s", string(e))
}

// DecodeBody checks the response Content-Type header.
// If there is no match, sets the error to ContentTypeError.
// Otherwise proceeds with Decode.
func (r *Response) DecodeBody(b BodyDecoderFrom) *Response {
	if r.Error != nil {
		return r
	}

	contentType := r.http.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, b.ContentType()) {
		r.Error = ContentTypeError(contentType)
		return r
	}

	return r.Decode(b)
}

func (r *Response) complete(err error) *Response {
	r.Error = err
	return r
}
