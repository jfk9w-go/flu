package flu

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// Client is a fluent http.Client wrapper.
type Client struct {
	*http.Client
	headers     http.Header
	statusCodes map[int]struct{}
}

// NewClient wraps the passed http.Client.
// If http == nil, creates a new http.Client
func NewClient(client *http.Client) *Client {
	if client == nil {
		client = &http.Client{Transport: NewTransport().http}
	}
	return &Client{
		Client:  client,
		headers: make(http.Header),
	}
}

// AddHeader allows to specify default headers set to every request.
func (c *Client) AddHeader(key, value string) *Client {
	c.headers.Add(key, value)
	return c
}

func (c *Client) AddHeaders(kvPairs ...string) *Client {
	l := keyValuePairsLength(kvPairs)
	for i := 0; i < l; i++ {
		c.AddHeader(kvPairs[2*i], kvPairs[2*i+1])
	}
	return c
}

func (c *Client) SetHeader(key, value string) *Client {
	c.headers.Set(key, value)
	return c
}

func (c *Client) SetHeaders(kvPairs ...string) *Client {
	l := keyValuePairsLength(kvPairs)
	for i := 0; i < l; i++ {
		c.SetHeader(kvPairs[2*i], kvPairs[2*i+1])
	}
	return c
}

func (c *Client) Timeout(timeout time.Duration) *Client {
	c.Client.Timeout = timeout
	return c
}

// SetCookies sets the http.Client cookies.
func (c *Client) SetCookies(rawurl string, cookies ...*http.Cookie) *Client {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	if c.Client.Jar == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			panic(err)
		}
		c.Client.Jar = jar
	}
	cookies = append(cookies, c.Client.Jar.Cookies(u)...)
	c.Client.Jar.SetCookies(u, cookies)
	return c
}

func (c *Client) AcceptResponseCodes(codes ...int) *Client {
	if c.statusCodes == nil {
		c.statusCodes = make(map[int]struct{})
	}
	for _, code := range codes {
		c.statusCodes[code] = struct{}{}
	}
	return c
}

func (c *Client) GET(resource string) *Request {
	return c.NewRequest(http.MethodGet, resource)
}

func (c *Client) HEAD(resource string) *Request {
	return c.NewRequest(http.MethodHead, resource)
}

func (c *Client) POST(resource string) *Request {
	return c.NewRequest(http.MethodPost, resource)
}

func (c *Client) PUT(resource string) *Request {
	return c.NewRequest(http.MethodPut, resource)
}

func (c *Client) PATCH(resource string) *Request {
	return c.NewRequest(http.MethodPatch, resource)
}

func (c *Client) DELETE(resource string) *Request {
	return c.NewRequest(http.MethodDelete, resource)
}

func (c *Client) CONNECT(resource string) *Request {
	return c.NewRequest(http.MethodConnect, resource)
}

func (c *Client) OPTIONS(resource string) *Request {
	return c.NewRequest(http.MethodOptions, resource)
}

func (c *Client) TRACE(resource string) *Request {
	return c.NewRequest(http.MethodTrace, resource)
}

// NewRequest creates a NewRequest.
func (c *Client) NewRequest(method string, resource string) *Request {
	return &Request{
		http:        c.Client,
		method:      method,
		resource:    resource,
		headers:     c.headers.Clone(),
		queryParams: make(url.Values),
		basicAuth:   [2]string{"", ""},
		statusCodes: c.statusCodes,
	}
}
