package flu

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

// Request allows to set basic http.Request properties.
type Request struct {
	client      *http.Client
	endpoint    string
	header      http.Header
	basicAuth   [2]string
	queryParams url.Values
	body        BodyWriter
	syncBody    bool
}

// Endpoint sets the request endpoint.
func (r *Request) Endpoint(endpoint string) *Request {
	r.endpoint = endpoint
	return r
}

// Header sets a request header.
func (r *Request) Header(key, value string) *Request {
	r.header.Set(key, value)
	return r
}

// BasicAuth allows to specify username and password to use in the basic authorization header.
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

// ReadBodyFunc sets the request body.
func (r *Request) Body(body BodyWriter) *Request {
	r.body = body
	return r
}

// Sync causes the request body to be loaded into a buffer before sending.
func (r *Request) Sync() *Request {
	r.syncBody = true
	return r
}

// Get executes a GET request.
func (r *Request) Get() *Response {
	return r.retrieve(http.MethodGet)
}

// Head executes a HEAD request.
func (r *Request) Head() *Response {
	return r.retrieve(http.MethodHead)
}

// Post executes a POST request.
func (r *Request) Post() *Response {
	return r.retrieve(http.MethodPost)
}

// Put executes a PUT request.
func (r *Request) Put() *Response {
	return r.retrieve(http.MethodPut)
}

// Patch executes a PATCH request.
func (r *Request) Patch() *Response {
	return r.retrieve(http.MethodPatch)
}

// Delete executes a DELETE request.
func (r *Request) Delete() *Response {
	return r.retrieve(http.MethodDelete)
}

// Connect executes a CONNECT request.
func (r *Request) Connect() *Response {
	return r.retrieve(http.MethodConnect)
}

// Options executes a OPTIONS request.
func (r *Request) Options() *Response {
	return r.retrieve(http.MethodOptions)
}

// Trace executes a TRACE request.
func (r *Request) Trace() *Response {
	return r.retrieve(http.MethodTrace)
}

func (r *Request) retrieve(method string) *Response {
	resp, err := r.exchange(method)
	return &Response{err, resp}
}

func (r *Request) exchange(method string) (*http.Response, error) {
	body, err := r.buildBody()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, r.endpoint, body)
	if err != nil {
		return nil, err
	}

	if r.body != nil {
		req.Header.Set("Content-Type", r.body.contentType())
	}

	if req.URL.RawQuery != "" {
		req.URL.RawQuery += "&"
	}

	req.URL.RawQuery += r.queryParams.Encode()
	if len(req.Header) == 0 {
		req.Header = r.header
	} else {
		for key, values := range r.header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	if r.basicAuth[0] != "" && r.basicAuth[1] != "" {
		req.SetBasicAuth(r.basicAuth[0], r.basicAuth[1])
	}

	return r.client.Do(req)
}

func (r *Request) buildBody() (io.Reader, error) {
	if r.body != nil {
		if r.syncBody {
			buf := new(bytes.Buffer)
			err := r.body.write(buf)
			if err != nil {
				return nil, err
			}

			return buf, nil
		}

		body, writer := io.Pipe()
		go func() {
			var err = r.body.write(writer)
			_ = writer.CloseWithError(err)
		}()

		return body, nil
	}

	return nil, nil
}
