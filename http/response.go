package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jfk9w-go/flu"
)

// Response is a fluent response wrapper.
type Response struct {
	*http.Response
	ignoreContentType bool

	// Error contains an error in case of a request processing error
	// or nil in case of success.
	Error error
}

func (r Response) IgnoreContentType() Response {
	if r.Error != nil {
		return r
	}
	r.ignoreContentType = true
	return r
}

type ResponseHandler interface {
	Handle(*http.Response) error
}

// HandleResponse executes a ResponseHandler if no previous handling errors occurred.
func (r Response) HandleResponse(handler ResponseHandler) Response {
	if r.Error != nil {
		return r
	}
	return r.complete(handler.Handle(r.Response))
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
		text := &flu.PlainText{""}
		if err := flu.DecodeFrom(flu.IO{R: r.Body}, text); err != nil {
			e.Text = fmt.Sprintf("response body read error: %s", err.Error())
		} else {
			e.Text = text.Value
		}
	}

	return e
}

// AcceptStatus checks the response status code and sets the error to StatusCodeError if there is no match.
func (r Response) AcceptStatus(codes ...int) Response {
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

// Decode reads the response body.
func (r Response) DecodeBody(decoder flu.DecoderFrom) Response {
	if r.Error != nil {
		return r
	}
	if !r.ignoreContentType {
		if c, ok := decoder.(ContentType); ok {
			contentType := r.Response.Header.Get("Content-Type")
			if !strings.HasPrefix(contentType, c.ContentType()) {
				return r.complete(ContentTypeError(c.ContentType()))
			}
		}
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

func (r Response) Reader() (io.Reader, error) {
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

func (r Response) DecodeBodyTo(out flu.Output) Response {
	if r.Error != nil {
		return r
	}
	return r.complete(flu.Copy(flu.IO{R: r.Body}, out))
}

func (r Response) complete(err error) Response {
	r.Error = err
	return r
}
