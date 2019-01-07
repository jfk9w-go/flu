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
	method      string
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

// Get sets the HTTP method to GET.
func (r *Request) Get() *Request {
	r.method = http.MethodGet
	return r
}

// Head sets the HTTP method to HEAD.
func (r *Request) Head() *Request {
	r.method = http.MethodHead
	return r
}

// Post sets the HTTP method to POST.
func (r *Request) Post() *Request {
	r.method = http.MethodPost
	return r
}

// Put sets the HTTP method to PUT.
func (r *Request) Put() *Request {
	r.method = http.MethodPut
	return r
}

// Patch sets the HTTP method to PATCH.
func (r *Request) Patch() *Request {
	r.method = http.MethodPatch
	return r
}

// Delete sets the HTTP method to DELETE.
func (r *Request) Delete() *Request {
	r.method = http.MethodDelete
	return r
}

// Connect sets the HTTP method to CONNECT.
func (r *Request) Connect() *Request {
	r.method = http.MethodConnect
	return r
}

// Options sets the HTTP method to OPTIONS.
func (r *Request) Options() *Request {
	r.method = http.MethodOptions
	return r
}

// Trace sets the HTTP method to TRACE.
func (r *Request) Trace() *Request {
	r.method = http.MethodTrace
	return r
}

// Execute executes the request and returns a response.
func (r *Request) Execute() *Response {
	resp, err := r.exchange()
	return &Response{err, resp}
}

func (r *Request) exchange() (*http.Response, error) {
	body, err := r.buildBody()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.method, r.endpoint, body)
	if err != nil {
		return nil, err
	}

	if r.body != nil {
		req.Header.Set("Content-Type", r.body.ContentType())
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
			err := r.body.Write(buf)
			if err != nil {
				return nil, err
			}

			return buf, nil
		}

		body, writer := io.Pipe()
		go func() {
			var err = r.body.Write(writer)
			_ = writer.CloseWithError(err)
		}()

		return body, nil
	}

	return nil, nil
}
