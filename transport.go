package flu

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Transport is a fluent wrapper around *http.Transport.
type Transport struct {
	httpTransport *http.Transport
}

// NewTransport initializes a new Transport with default settings.
// This should be equivalent to http.DefaultTransport
func NewTransport() *Transport {
	httpTransport := new(http.Transport)
	*httpTransport = *http.DefaultTransport.(*http.Transport)
	return &Transport{httpTransport}
}

func (t *Transport) RoundTrip(httpReq *http.Request) (*http.Response, error) {
	return t.httpTransport.RoundTrip(httpReq)
}

// Proxy sets the http.Transport.Proxy.
func (t *Transport) Proxy(proxy func(*http.Request) (*url.URL, error)) *Transport {
	t.httpTransport.Proxy = proxy
	return t
}

// DialContext sets http.Transport.DialContext.
func (t *Transport) DialContext(dialContext func(ctx context.Context, network, addr string) (net.Conn, error)) *Transport {
	t.httpTransport.DialContext = dialContext
	return t
}

// MaxIdleConns sets http.Transport.MaxIdleConns.
func (t *Transport) MaxIdleConns(maxIdleConns int) *Transport {
	t.httpTransport.MaxIdleConns = maxIdleConns
	return t
}

// MaxIdleConnsPerHost sets http.Transport.MaxIdleConnsPerHost.
func (t *Transport) MaxIdleConnsPerHost(maxIdleConnsPerHost int) *Transport {
	t.httpTransport.MaxConnsPerHost = maxIdleConnsPerHost
	return t
}

// MaxConnsPerHost sets http.Transport.MaxConnsPerHost.
func (t *Transport) MaxConnsPerHost(maxConnsPerHost int) *Transport {
	t.httpTransport.MaxConnsPerHost = maxConnsPerHost
	return t
}

// IdleConnTimeout sets http.Transport.IdleConnTimeout.
func (t *Transport) IdleConnTimeout(idleConnTimeout time.Duration) *Transport {
	t.httpTransport.IdleConnTimeout = idleConnTimeout
	return t
}

// ResponseHeaderTimeout sets http.Transport.ResponseHeaderTimeout.
func (t *Transport) ResponseHeaderTimeout(responseHeaderTimeout time.Duration) *Transport {
	t.httpTransport.ResponseHeaderTimeout = responseHeaderTimeout
	return t
}

// TLSHandshakeTimeout sets http.Transport.TLSHandshakeTimeout.
func (t *Transport) TLSHandshakeTimeout(tlsHandshakeTimeout time.Duration) *Transport {
	t.httpTransport.TLSHandshakeTimeout = tlsHandshakeTimeout
	return t
}

// ExpectContinueTimeout sets http.Transport.ExpectContinueTimeout.
func (t *Transport) ExpectContinueTimeout(expectContinueTimeout time.Duration) *Transport {
	t.httpTransport.ExpectContinueTimeout = expectContinueTimeout
	return t
}

// NewClient creates a new Client with this Transport.
func (t *Transport) NewClient() *Client {
	return NewClient(&http.Client{Transport: t.httpTransport})
}
