package flu

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

// Request is a fluent http.Request wrapper.
type Request interface {
	// Retrieve executes the request and returns a Response.
	Retrieve() Response
}

// BaseRequest allows to set basic http.Request properties.
type BaseRequest struct {
	endpoint    string
	header      http.Header
	basicAuth   [2]string
	queryParams url.Values
	client      *http.Client
}

// Endpoint sets the request endpoint.
func (base BaseRequest) Endpoint(endpoint string) BaseRequest {
	base.endpoint = endpoint
	return base
}

// Header sets a request header.
func (base BaseRequest) Header(key, value string) BaseRequest {
	base.header.Set(key, value)
	return base
}

// BasicAuth allows to specify username and password to use in the basic authorization header.
func (base BaseRequest) BasicAuth(username, password string) BaseRequest {
	base.basicAuth[0] = username
	base.basicAuth[1] = password
	return base
}

// QueryParam sets a query parameter.
func (base BaseRequest) QueryParam(key, value string) BaseRequest {
	base.queryParams.Add(key, value)
	return base
}

func (base BaseRequest) exchange(req *http.Request) (*http.Response, error) {
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

// Get returns a GET request builder.
func (base BaseRequest) Get() GetRequest {
	return GetRequest(base)
}

// Post returns a POST request builder.
func (base BaseRequest) Post() PostRequest {
	return PostRequest{base: base}
}

// GetRequest is a fluent GET http.Request wrapper.
type GetRequest BaseRequest

func (get GetRequest) exchange() (*http.Response, error) {
	var req, err = http.NewRequest(http.MethodGet, get.endpoint, nil)
	if err != nil {
		return nil, err
	}

	return BaseRequest(get).exchange(req)
}

func (get GetRequest) Retrieve() Response {
	var resp, err = get.exchange()
	return Response{resp, err}
}

// PostRequest is a fluent POST http.Request wrapper.
type PostRequest struct {
	base BaseRequest
	body RequestBodyBuilder
	sync bool
}

// Body sets the request body.
func (post PostRequest) Body(body RequestBodyBuilder) PostRequest {
	post.body = body
	return post
}

func (post PostRequest) Sync() PostRequest {
	post.sync = true
	return post
}

func (post PostRequest) exchange() (resp *http.Response, err error) {
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

func (post PostRequest) buildBody() (body io.Reader, err error) {
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

func (post PostRequest) Retrieve() Response {
	var resp, err = post.exchange()
	return Response{resp, err}
}
