package flu

import (
	"bytes"
	"context"
	"io"
	"log"
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
	body        BodyEncoderTo
	useBuffer   bool
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

// AddHeaders adds request headers.
// keyValues is an array of key-value pairs and must have even length.
func (r *Request) AddHeaders(keyValues ...string) *Request {
	keyValuesLength := len(keyValues)
	if keyValuesLength%2 == 1 {
		log.Fatal("keyValues length must be even, got ", keyValuesLength)
	}

	for i := 0; i < keyValuesLength; i += 2 {
		key, value := keyValues[i], keyValues[i+1]
		r.AddHeader(key, value)
	}

	return r
}

// SetHeaders sets request headers.
// keyValues is an array of key-value pairs and must have even length.
func (r *Request) SetHeaders(keyValues ...string) *Request {
	keyValuesLength := len(keyValues)
	if keyValuesLength%2 == 1 {
		log.Fatal("keyValues length must be even, got ", keyValuesLength)
	}

	for i := 0; i < keyValuesLength; i += 2 {
		key, value := keyValues[i], keyValues[i+1]
		r.SetHeader(key, value)
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
	r.body = body
	return r
}

// Buffer causes the request body to be loaded into a buffer before sending.
func (r *Request) Buffer() *Request {
	r.useBuffer = true
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
	body, err := r.buildBody()
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

	return r.http.Do(httpReq)
}

func (r *Request) buildBody() (io.Reader, error) {
	if r.body != nil {
		if r.useBuffer {
			buf := new(bytes.Buffer)
			err := r.body.EncodeTo(buf)
			if err != nil {
				return nil, err
			}

			return buf, nil
		}

		body, writer := io.Pipe()
		go func() {
			err := r.body.EncodeTo(writer)
			_ = writer.CloseWithError(err)
		}()

		return body, nil
	}

	return nil, nil
}
