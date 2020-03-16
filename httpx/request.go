package httpx

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/jfk9w-go/flu"
)

// NewRequest allows to set basic http.NewRequest properties.
type Request struct {
	*http.Request
	client Client
	query  url.Values
	body   flu.Body
	err    error
}

// AddHeader adds a request header.
func (r Request) AddHeader(key, value string) Request {
	if r.err != nil {
		return r
	}
	r.Header.Add(key, value)
	return r
}

// SetHeader sets a request header.
func (r Request) SetHeader(key, value string) Request {
	if r.err != nil {
		return r
	}
	r.Header.Set(key, value)
	return r
}

type VarargsLengthError int

func (e VarargsLengthError) Error() string {
	return fmt.Sprintf("key-value pairs array length must be even, got %d", e.Length())
}

func (e VarargsLengthError) Length() int {
	return int(e)
}

func keyValuePairsLength(kvPairs []string) int {
	length := len(kvPairs)
	if length%2 > 0 {
		panic(VarargsLengthError(length))
	}
	return length
}

// AddHeaders adds request header.
// kvPairs is an array of key-value pairs and must have even length.
func (r Request) AddHeaders(kvPairs ...string) Request {
	if r.err != nil {
		return r
	}
	kvLength := keyValuePairsLength(kvPairs)
	for i := 0; i < kvLength; i += 2 {
		k, v := kvPairs[i], kvPairs[i+1]
		r.AddHeader(k, v)
	}
	return r
}

// SetHeaders sets request header.
// kvPairs is an array of key-value pairs and must have even length.
func (r Request) SetHeaders(kvPairs ...string) Request {
	if r.err != nil {
		return r
	}
	kvLength := keyValuePairsLength(kvPairs)
	for i := 0; i < kvLength; i += 2 {
		k, v := kvPairs[i], kvPairs[i+1]
		r.SetHeader(k, v)
	}
	return r
}

func (r Request) Auth(auth Authorization) Request {
	if r.err != nil {
		return r
	}
	auth.SetAuth(r.Request)
	return r
}

// QueryParam sets a query parameter.
func (r Request) QueryParam(key, value string) Request {
	if r.err != nil {
		return r
	}
	r.query.Set(key, value)
	return r
}

// Body sets the request body.
func (r Request) Body(body flu.Body) Request {
	if r.err != nil {
		return r
	}
	switch body.(type) {
	case flu.BodyEncoderTo:
	case flu.Readable:
	default:
		panic(errors.Errorf("unrecognized body type: %T", body))
	}
	r.body = body
	return r
}

func (r Request) Context(ctx context.Context) Request {
	if r.err != nil {
		return r
	}
	r.Request = r.Request.WithContext(ctx)
	return r
}

// Send executes the request and returns a response.
func (r Request) Execute() Response {
	resp, err := r.do()
	return Response{resp, err}
}

func (r Request) do() (*http.Response, error) {
	if r.err != nil {
		return nil, r.err
	}

	if r.body != nil {
		if cl, ok := r.body.(ContentLength); ok {
			contentLength, err := cl.ContentLength()
			if err != nil {
				return nil, err
			}
			r.Request.ContentLength = contentLength
		}

		if b, ok := r.body.(flu.Readable); ok {
			body, err := b.Reader()
			if err != nil {
				return nil, err
			}
			if body, ok := body.(io.ReadCloser); ok {
				r.Request.Body = body
			} else {
				r.Request.Body = ioutil.NopCloser(body)
			}
			if sized, ok := body.(interface{ Len() int }); ok {
				r.Request.ContentLength = int64(sized.Len())
			}
		} else if b, ok := r.body.(flu.BodyEncoderTo); ok {
			body, err := flu.ReadablePipe(b).Reader()
			if err != nil {
				return nil, err
			}
			r.Request.Body = body.(io.ReadCloser)
		}

		r.Request.Header.Set("Content-Type", r.body.ContentType())
	}

	r.Request.URL.RawQuery = r.query.Encode()
	response, err := r.client.Do(r.Request)
	if err != nil {
		return nil, err
	}
	if len(r.client.statuses) > 0 {
		if !r.client.statuses[response.StatusCode] {
			return nil, NewStatusCodeError(response)
		}
	}

	return response, nil
}

type ContentLength interface {
	ContentLength() (int64, error)
}
