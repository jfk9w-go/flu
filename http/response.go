package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/jfk9w-go/flu"
)

// Response is a fluent response wrapper.
type Response struct {
	*http.Response

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
	return r.complete(handler.Handle(r.Response))
}

type StatusCodeError struct {
	StatusCode   int
	ResponseBody flu.Bytes
}

func (e StatusCodeError) Error() string {
	text := fmt.Sprintf("%d", e.StatusCode)
	if len(e.ResponseBody) > 0 {
		body := string(e.ResponseBody)
		text += " (" + body + ")"
	}

	return text
}

func NewStatusCodeError(r *http.Response) StatusCodeError {
	e := StatusCodeError{StatusCode: r.StatusCode}
	if r.Body != nil {
		body := new(flu.ByteBuffer)
		if _, err := flu.Copy(flu.IO{R: r.Body}, body); err != nil {
			log.Printf("Failed to read response body from %s: %s", r.Request.URL, err)
		} else {
			e.ResponseBody = body.Bytes()
		}
	}

	return e
}

// CheckStatus checks the response status code and sets the error to StatusCodeError if there is no match.
func (r *Response) CheckStatus(codes ...int) *Response {
	if r.Error != nil {
		return r
	}
	for _, c := range codes {
		if c == r.StatusCode {
			return r
		}
	}

	return r.complete(NewStatusCodeError(r.Response))
}

func (r *Response) CheckContentType(value string) *Response {
	if r.Error != nil {
		return r
	}
	contentType := r.Response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, value) {
		return r.complete(ContentTypeError(contentType))
	}
	return r
}

// Decode reads the response body.
func (r *Response) DecodeBody(decoder flu.DecoderFrom) *Response {
	if r.Error != nil {
		return r
	}
	return r.complete(flu.DecodeFrom(flu.IO{R: r.Body}, decoder))
}

type ContentTypeError string

func (e ContentTypeError) Error() string {
	return fmt.Sprintf("invalid body type: %s", string(e))
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

func (r *Response) DecodeBodyTo(out flu.Output) *Response {
	if r.Error != nil {
		return r
	}
	_, err := flu.Copy(flu.IO{R: r.Body}, out)
	return r.complete(err)
}

func (r *Response) complete(err error) *Response {
	r.Error = err
	return r
}
