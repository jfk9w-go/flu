package flu

import (
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// Client is a fluent http.Client wrapper.
type Client struct {
	client         *http.Client
	defaultHeaders http.Header
}

// NewClient wraps the passed http.Client.
// If client == nil, creates a new http.Client
func NewClient(client *http.Client) *Client {
	if client == nil {
		client = new(http.Client)
	}

	return &Client{
		client:         client,
		defaultHeaders: http.Header{},
	}
}

// Cookies sets the client cookies.
func (c *Client) Cookies(rawurl string, cookies ...*http.Cookie) *Client {
	url, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	if c.client.Jar == nil {
		var jar, err = cookiejar.New(nil)
		if err != nil {
			panic(err)
		}

		c.client.Jar = jar
	}

	cookies = append(cookies, c.client.Jar.Cookies(url)...)
	c.client.Jar.SetCookies(url, cookies)

	return c
}

// ResponseHeaderTimeout allows to specify the response header timeout on the client.
func (c *Client) ResponseHeaderTimeout(timeout time.Duration) *Client {
	if timeout > 0 {
		c.transport().ResponseHeaderTimeout = timeout
	}

	return c
}

// DefaultHeader allows to specify a default header set to every request.
func (c *Client) DefaultHeader(key, value string) *Client {
	c.defaultHeaders.Set(key, value)
	return c
}

// NewRequest creates a Request builder.
func (c *Client) NewRequest() *BaseRequest {
	base := &BaseRequest{
		client:      c.client,
		header:      http.Header{},
		queryParams: url.Values{},
		basicAuth:   [2]string{"", ""},
	}

	for key, values := range c.defaultHeaders {
		for _, value := range values {
			base.Header(key, value)
		}
	}

	return base
}

func (c *Client) transport() *http.Transport {
	if c.client.Transport == nil {
		var transport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: time.Second,
			ResponseHeaderTimeout: time.Minute,
		}

		c.client.Transport = transport
		return transport
	} else if transport, ok := c.client.Transport.(*http.Transport); ok {
		return transport
	} else {
		panic(fmt.Errorf("invalid transport type: %T", c.client.Transport))
	}
}
