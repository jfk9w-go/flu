package flu

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

// BaseRequest allows to set basic http.Request properties.
type BaseRequest struct {
	client      *http.Client
	endpoint    string
	header      http.Header
	basicAuth   [2]string
	queryParams url.Values
}

// Endpoint sets the request endpoint.
func (base *BaseRequest) Endpoint(endpoint string) *BaseRequest {
	base.endpoint = endpoint
	return base
}

// Header sets a request header.
func (base *BaseRequest) Header(key, value string) *BaseRequest {
	base.header.Set(key, value)
	return base
}

// BasicAuth allows to specify username and password to use in the basic authorization header.
func (base *BaseRequest) BasicAuth(username, password string) *BaseRequest {
	base.basicAuth[0] = username
	base.basicAuth[1] = password
	return base
}

// QueryParam sets a query parameter.
func (base *BaseRequest) QueryParam(key, value string) *BaseRequest {
	base.queryParams.Add(key, value)
	return base
}

func (base *BaseRequest) exchange(req *http.Request) (*http.Response, error) {
	if req.URL.RawQuery != "" {
		req.URL.RawQuery += "&"
	}

	req.URL.RawQuery += base.queryParams.Encode()
	if len(req.Header) == 0 {
		req.Header = base.header
	} else {
		for key, values := range base.header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	if base.basicAuth[0] != "" && base.basicAuth[1] != "" {
		req.SetBasicAuth(base.basicAuth[0], base.basicAuth[1])
	}

	return base.client.Do(req)
}

// GET returns a GET request builder.
func (base *BaseRequest) GET() *GET {
	return (*GET)(base)
}

// POST returns a POST request builder.
func (base *BaseRequest) POST() *POST {
	return &POST{base: base}
}

// GET is a fluent GET http.Request wrapper.
type GET BaseRequest

// Retrieve sends the request and waits for a response.
func (get *GET) Retrieve() *Response {
	resp, err := get.exchange()
	return &Response{err, resp}
}

func (get *GET) exchange() (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, get.endpoint, nil)
	if err != nil {
		return nil, err
	}

	return (*BaseRequest)(get).exchange(req)
}

// POST is a fluent POST http.Request wrapper.
type POST struct {
	base *BaseRequest
	body RequestBodyBuilder
	sync bool
}

// Body sets the request body.
func (post *POST) Body(body RequestBodyBuilder) *POST {
	post.body = body
	return post
}

// Sync causes the body to be loaded into a buffer first.
func (post *POST) Sync() *POST {
	post.sync = true
	return post
}

// Retrieve sends the request and waits for a response.
func (post *POST) Retrieve() *Response {
	resp, err := post.exchange()
	return &Response{err, resp}
}

func (post *POST) exchange() (resp *http.Response, err error) {
	var body io.Reader
	body, err = post.buildBody()
	if err != nil {
		return
	}

	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, post.base.endpoint, body)
	if err != nil {
		return
	}

	if post.body != nil {
		req.Header.Set("Content-Type", post.body.contentType())
	}

	return post.base.exchange(req)
}

func (post *POST) buildBody() (body io.Reader, err error) {
	if post.body != nil {
		if post.sync {
			var buf = new(bytes.Buffer)
			err = post.body.build(buf)
			if err != nil {
				return
			}

			body = buf
			return
		}

		var writer *io.PipeWriter
		body, writer = io.Pipe()
		go func() {
			var err = post.body.build(writer)
			_ = writer.CloseWithError(err)
		}()
	}

	return
}
