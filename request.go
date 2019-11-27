package flu

import (
	"bytes"
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
	body        BodyEncoderTo
	buffer      bool
	respCodes   map[int]struct{}
}

// Resource sets the request resource.
func (req *Request) Resource(resource string) *Request {
	req.resource = resource
	return req
}

// AddHeader adds a request header.
func (req *Request) AddHeader(key, value string) *Request {
	req.headers.Add(key, value)
	return req
}

// SetHeader sets a request header.
func (req *Request) SetHeader(key, value string) *Request {
	req.headers.Set(key, value)
	return req
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
func (req *Request) AddHeaders(kvPairs ...string) *Request {
	kvLength := keyValuePairsLength(kvPairs)
	for i := 0; i < kvLength; i += 2 {
		k, v := kvPairs[i], kvPairs[i+1]
		req.AddHeader(k, v)
	}

	return req
}

// SetHeaders sets request headers.
// kvPairs is an array of key-value pairs and must have even length.
func (req *Request) SetHeaders(kvPairs ...string) *Request {
	kvLength := keyValuePairsLength(kvPairs)
	for i := 0; i < kvLength; i += 2 {
		k, v := kvPairs[i], kvPairs[i+1]
		req.SetHeader(k, v)
	}

	return req
}

// BasicAuth allows to specify username and password to use in the basic authorization headers.
func (req *Request) BasicAuth(username, password string) *Request {
	req.basicAuth[0] = username
	req.basicAuth[1] = password
	return req
}

// QueryParam sets a query parameter.
func (req *Request) QueryParam(key, value string) *Request {
	req.queryParams.Add(key, value)
	return req
}

// Body sets the request body.
func (req *Request) Body(body BodyEncoderTo) *Request {
	req.body = body
	return req
}

// Buffer causes the request body to be loaded into a buffer before sending.
func (req *Request) Buffer() *Request {
	req.buffer = true
	return req
}

// GET sets the HTTP method to GET.
func (req *Request) GET() *Request {
	req.method = http.MethodGet
	return req
}

// HEAD sets the HTTP method to HEAD.
func (req *Request) HEAD() *Request {
	req.method = http.MethodHead
	return req
}

// POST sets the HTTP method to POST.
func (req *Request) POST() *Request {
	req.method = http.MethodPost
	return req
}

// PUT sets the HTTP method to PUT.
func (req *Request) PUT() *Request {
	req.method = http.MethodPut
	return req
}

// PATCH sets the HTTP method to PATCH.
func (req *Request) PATCH() *Request {
	req.method = http.MethodPatch
	return req
}

// DELETE sets the HTTP method to DELETE.
func (req *Request) DELETE() *Request {
	req.method = http.MethodDelete
	return req
}

// CONNECT sets the HTTP method to CONNECT.
func (req *Request) CONNECT() *Request {
	req.method = http.MethodConnect
	return req
}

// OPTIONS sets the HTTP method to OPTIONS.
func (req *Request) OPTIONS() *Request {
	req.method = http.MethodOptions
	return req
}

// TRACE sets the HTTP method to TRACE.
func (req *Request) TRACE() *Request {
	req.method = http.MethodTrace
	return req
}

// Send executes the request and returns a response.
func (req *Request) Send() *Response {
	resp, err := req.send(nil)
	return &Response{resp, err}
}

func (req *Request) SendWithContext(ctx context.Context) *Response {
	httpResp, err := req.send(ctx)
	return &Response{httpResp, err}
}

func (req *Request) send(ctx context.Context) (*http.Response, error) {
	body, err := req.content()
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(req.method, req.resource, body)
	if err != nil {
		return nil, err
	}

	if req.body != nil {
		httpReq.Header.Set("Content-Type", req.body.ContentType())
		if httpReq.ContentLength == 0 {
			// see net/http/request.go:896
			httpReq.ContentLength = -1
		}
	}

	if httpReq.URL.RawQuery != "" {
		httpReq.URL.RawQuery += "&"
	}

	httpReq.URL.RawQuery += req.queryParams.Encode()
	if len(httpReq.Header) == 0 {
		httpReq.Header = req.headers
	} else {
		for key, values := range req.headers {
			for _, value := range values {
				httpReq.Header.Add(key, value)
			}
		}
	}

	if req.basicAuth[0] != "" && req.basicAuth[1] != "" {
		httpReq.SetBasicAuth(req.basicAuth[0], req.basicAuth[1])
	}

	if ctx != nil {
		httpReq = httpReq.WithContext(ctx)
	}

	httpResp, err := req.http.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if req.respCodes != nil {
		if _, ok := req.respCodes[httpResp.StatusCode]; !ok {
			return nil, createStatusCodeError(httpResp)
		}
	}

	return httpResp, nil
}

func (req *Request) content() (io.Reader, error) {
	if req.body != nil {
		if req.buffer {
			buf := new(bytes.Buffer)
			err := req.body.EncodeTo(buf)
			if err != nil {
				return nil, err
			}

			return buf, nil
		}

		body, writer := io.Pipe()
		go func() {
			err := req.body.EncodeTo(writer)
			_ = writer.CloseWithError(err)
		}()

		return body, nil
	}

	return nil, nil
}
