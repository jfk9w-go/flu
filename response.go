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
func (resp *Response) HandleResponse(handler ResponseHandler) *Response {
	if resp.Error != nil {
		return resp
	}

	return resp.complete(handler.Handle(resp.http))
}

type StatusCodeError struct {
	Code int
	Text string
}

func (e StatusCodeError) Error() string {
	text := fmt.Sprintf("invalid status code: %d", int(e.Code))
	if e.Text != "" {
		text += " (" + e.Text + ")"
	}

	return text
}

func createStatusCodeError(resp *http.Response) StatusCodeError {
	e := StatusCodeError{Code: resp.StatusCode}
	if resp.Body != nil {
		defer resp.Body.Close()
		text := PlainText("")
		err := text.DecodeFrom(resp.Body)
		if err != nil {
			text.Value = fmt.Sprintf("response body read error: %s", err.Error())
		}

		e.Text = text.Value
	}

	return e
}

// CheckStatusCode checks the response status code and sets the error to StatusCodeError if there is no match.
func (resp *Response) CheckStatusCode(okCodes ...int) *Response {
	if resp.Error != nil {
		return resp
	}

	for _, c := range okCodes {
		if c == resp.http.StatusCode {
			return resp
		}
	}

	return resp.complete(createStatusCodeError(resp.http))
}

// Decode decodes the response body.
func (resp *Response) Decode(decoder DecoderFrom) *Response {
	if resp.Error != nil {
		return resp
	}

	body := resp.http.Body
	defer body.Close()
	return resp.complete(decoder.DecodeFrom(body))
}

type ContentTypeError string

func (e ContentTypeError) Error() string {
	return fmt.Sprintf("invalid content type: %s", string(e))
}

// DecodeBody checks the response Content-Type header.
// If there is no match, sets the error to ContentTypeError.
// Otherwise proceeds with Decode.
func (resp *Response) DecodeBody(body BodyDecoderFrom) *Response {
	if resp.Error != nil {
		return resp
	}

	contentType := resp.http.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, body.ContentType()) {
		resp.Error = ContentTypeError(contentType)
		return resp
	}

	return resp.Decode(body)
}

type WriteResourceError struct {
	Err error
}

func (e WriteResourceError) Error() string {
	return fmt.Sprintf("failed to read resource from body: %s", e.Err)
}

func (resp *Response) ReadResource(res ResourceWriter) *Response {
	if resp.Error != nil {
		return resp
	}

	w, err := res.Writer()
	if err != nil {
		return resp.complete(WriteResourceError{err})
	}

	//noinspection GoUnhandledErrorResult
	defer w.Close()

	body := resp.http.Body
	//noinspection GoUnhandledErrorResult
	defer body.Close()

	_, err = io.Copy(w, body)
	return resp.complete(err)
}

func (resp *Response) complete(err error) *Response {
	resp.Error = err
	return resp
}
