package flu

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Request allows to set basic http.Request properties.
type Request struct {
	http        *http.Client
	method      string
	resource    string
	headers     http.Header
	basicAuth   [2]string
	queryParams url.Values
	body        BodyWriter
	buffer      bool
	statusCodes map[int]struct{}
}

// Resource sets the request resource.
func (r *Request) Resource(resource string) *Request {
	r.resource = resource
	return r
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
func (r *Request) Body(body BodyWriter) *Request {
	r.body = body
	return r
}

// Buffer causes the request body to be loaded into a buffer before sending.
func (r *Request) Buffer() *Request {
	r.buffer = true
	return r
}

// GET sets the HTTP method to GET.
func (r *Request) GET() *Request {
	r.method = http.MethodGet
	return r
}

// HEAD sets the HTTP method to HEAD.
func (r *Request) HEAD() *Request {
	r.method = http.MethodHead
	return r
}

// POST sets the HTTP method to POST.
func (r *Request) POST() *Request {
	r.method = http.MethodPost
	return r
}

// PUT sets the HTTP method to PUT.
func (r *Request) PUT() *Request {
	r.method = http.MethodPut
	return r
}

// PATCH sets the HTTP method to PATCH.
func (r *Request) PATCH() *Request {
	r.method = http.MethodPatch
	return r
}

// DELETE sets the HTTP method to DELETE.
func (r *Request) DELETE() *Request {
	r.method = http.MethodDelete
	return r
}

// CONNECT sets the HTTP method to CONNECT.
func (r *Request) CONNECT() *Request {
	r.method = http.MethodConnect
	return r
}

// OPTIONS sets the HTTP method to OPTIONS.
func (r *Request) OPTIONS() *Request {
	r.method = http.MethodOptions
	return r
}

// TRACE sets the HTTP method to TRACE.
func (r *Request) TRACE() *Request {
	r.method = http.MethodTrace
	return r
}

// Send executes the request and returns a response.
func (r *Request) Send() *Response {
	resp, err := r.send(nil)
	return &Response{resp, err}
}

func (r *Request) SendWithContext(ctx context.Context) *Response {
	httpResp, err := r.send(ctx)
	return &Response{httpResp, err}
}

func (r *Request) send(ctx context.Context) (*http.Response, error) {
	body, err := r.content()
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest(r.method, r.resource, body)
	if err != nil {
		return nil, err
	}
	if r.body != nil {
		httpReq.Header.Set("Content-Type", r.body.ContentType())
	}
	if httpReq.URL.RawQuery != "" {
		httpReq.URL.RawQuery += "&"
	}
	httpReq.URL.RawQuery += r.queryParams.Encode()
	if len(httpReq.Header) == 0 {
		httpReq.Header = r.headers
	} else {
		for key, values := range r.headers {
			for _, value := range values {
				httpReq.Header.Add(key, value)
			}
		}
	}
	if r.basicAuth[0] != "" && r.basicAuth[1] != "" {
		httpReq.SetBasicAuth(r.basicAuth[0], r.basicAuth[1])
	}
	if ctx != nil {
		httpReq = httpReq.WithContext(ctx)
	}
	httpResp, err := r.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if r.statusCodes != nil {
		if _, ok := r.statusCodes[httpResp.StatusCode]; !ok {
			return nil, createStatusCodeError(httpResp)
		}
	}
	return httpResp, nil
}

func (r *Request) content() (io.Reader, error) {
	if r.body == nil {
		return nil, nil
	}
	if !r.buffer {
		return PipeOut(r.body).Reader()
	}
	buf := new(Buffer)
	err := Write(r.body, buf)
	if err != nil {
		return nil, err
	}
	return buf.Reader()
}
