package flu

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// NewRequest allows to set basic http.NewRequest properties.
type Request struct {
	http        *http.Client
	method      string
	resource    string
	headers     http.Header
	basicAuth   [2]string
	queryParams url.Values
	bodyEncoder BodyEncoderTo
	buffer      bool
	statusCodes map[int]struct{}
}

// AddHeader adds a request header.
func (r *Request) AddHeader(key, value string) *Request {
	r.headers.Add(key, value)
	return r
}

// SetHeader sets a request header.
func (r *Request) SetHeader(key, value string) *Request {
	r.headers.Set(key, value)
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

// AddHeaders adds request headers.
// kvPairs is an array of key-value pairs and must have even length.
func (r *Request) AddHeaders(kvPairs ...string) *Request {
	kvLength := keyValuePairsLength(kvPairs)
	for i := 0; i < kvLength; i += 2 {
		k, v := kvPairs[i], kvPairs[i+1]
		r.AddHeader(k, v)
	}
	return r
}

// SetHeaders sets request headers.
// kvPairs is an array of key-value pairs and must have even length.
func (r *Request) SetHeaders(kvPairs ...string) *Request {
	kvLength := keyValuePairsLength(kvPairs)
	for i := 0; i < kvLength; i += 2 {
		k, v := kvPairs[i], kvPairs[i+1]
		r.SetHeader(k, v)
	}
	return r
}

// BasicAuth allows to specify username and password to use in the basic authorization headers.
func (r *Request) BasicAuth(username, password string) *Request {
	r.basicAuth[0] = username
	r.basicAuth[1] = password
	return r
}

// QueryParam sets a query parameter.
func (r *Request) QueryParam(key, value string) *Request {
	r.queryParams.Add(key, value)
	return r
}

// Body sets the request body.
func (r *Request) Body(body BodyEncoderTo) *Request {
	r.bodyEncoder = body
	return r
}

// Buffer causes the request body to be loaded into a buffer before sending.
func (r *Request) Buffer() *Request {
	r.buffer = true
	return r
}

// Send executes the request and returns a response.
func (r *Request) Execute() *Response {
	resp, err := r.do(nil)
	return &Response{resp, err}
}

func (r *Request) SendWithContext(ctx context.Context) *Response {
	resp, err := r.do(ctx)
	return &Response{resp, err}
}

func (r *Request) do(ctx context.Context) (*http.Response, error) {
	body, err := r.content()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(r.method, r.resource, body)
	if err != nil {
		return nil, err
	}
	if r.bodyEncoder != nil {
		req.Header.Set("Content-Type", r.bodyEncoder.ContentType())
	}
	if req.URL.RawQuery != "" {
		req.URL.RawQuery += "&"
	}
	req.URL.RawQuery += r.queryParams.Encode()
	if len(req.Header) == 0 {
		req.Header = r.headers
	} else {
		for key, values := range r.headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}
	if r.basicAuth[0] != "" && r.basicAuth[1] != "" {
		req.SetBasicAuth(r.basicAuth[0], r.basicAuth[1])
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	resp, err := r.http.Do(req)
	if err != nil {
		return nil, err
	}
	if r.statusCodes != nil {
		if _, ok := r.statusCodes[resp.StatusCode]; !ok {
			return nil, NewStatusCodeError(resp)
		}
	}
	return resp, nil
}

func (r *Request) content() (io.Reader, error) {
	if r.bodyEncoder == nil {
		return nil, nil
	}
	var body Readable
	if r.buffer {
		if bb, ok := r.bodyEncoder.(bufferedBody); ok {
			body = bb.buf
		} else {
			buf := NewBuffer()
			err := EncodeTo(r.bodyEncoder, buf)
			if err != nil {
				return nil, err
			}
			r.bodyEncoder = bufferedBody{buf, r.bodyEncoder.ContentType()}
			body = buf
		}
	} else {
		body = Input(r.bodyEncoder)
	}
	return body.Reader()
}

type bufferedBody struct {
	buf         Buffer
	contentType string
}

func (bb bufferedBody) EncodeTo(w io.Writer) error {
	_, err := bb.buf.WriteTo(w)
	return err
}

func (bb bufferedBody) ContentType() string {
	return bb.contentType
}
