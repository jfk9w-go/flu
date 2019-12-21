package flu

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// Client is a fluent http.Client wrapper.
type Client struct {
	http        *http.Client
	headers     []string
	statusCodes map[int]struct{}
}

// NewClient wraps the passed http.Client.
// If http == nil, creates a new http.Client
func NewClient(client *http.Client) *Client {
	if client == nil {
		client = &http.Client{Transport: NewTransport().http}
	}
	return &Client{
		http:    client,
		headers: make([]string, 0),
	}
}

// AddHeader allows to specify default headers set to every request.
func (c *Client) AddHeader(key, value string) *Client {
	c.headers = append(c.headers, key, value)
	return c
}

func (c *Client) AddHeaders(kvPairs ...string) *Client {
	keyValuePairsLength(kvPairs)
	c.headers = append(c.headers, kvPairs...)
	return c
}

func (c *Client) Timeout(timeout time.Duration) *Client {
	c.http.Timeout = timeout
	return c
}

// SetCookies sets the http.Client cookies.
func (c *Client) SetCookies(rawurl string, cookies ...*http.Cookie) *Client {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	if c.http.Jar == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			panic(err)
		}
		c.http.Jar = jar
	}
	cookies = append(cookies, c.http.Jar.Cookies(u)...)
	c.http.Jar.SetCookies(u, cookies)
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

// NewRequest creates a Request.
func (c *Client) NewRequest() *Request {
	req := &Request{
		http:        c.http,
		method:      http.MethodGet,
		headers:     http.Header{},
		queryParams: url.Values{},
		basicAuth:   [2]string{"", ""},
		statusCodes: c.statusCodes,
	}
	req.SetHeaders(c.headers...)
	return req
}
