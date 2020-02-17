package flu

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Transport is a fluent wrapper around *http.Transport.
type Transport struct {
	http    *http.Transport
	limiter Limiter
}

// NewTransport initializes a new Transport with default settings.
// This should be equivalent to http.DefaultTransport
func NewTransport() *Transport {
	return &Transport{
		http:    http.DefaultTransport.(*http.Transport).Clone(),
		limiter: Unlimiter,
	}
}

// Proxy sets the http.Transport.Proxy.
func (t *Transport) Proxy(proxy func(*http.Request) (*url.URL, error)) *Transport {
	t.http.Proxy = proxy
	return t
}

func (t *Transport) ProxyURL(rawurl string) *Transport {
	if rawurl == "" {
		return t.Proxy(nil)
	}
	proxy, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return t.Proxy(http.ProxyURL(proxy))
}

// DialContext sets http.Transport.DialContext.
func (t *Transport) DialContext(fun func(ctx context.Context, network, addr string) (net.Conn, error)) *Transport {
	t.http.DialContext = fun
	return t
}

// MaxIdleConns sets http.Transport.MaxIdleConns.
func (t *Transport) MaxIdleConns(value int) *Transport {
	t.http.MaxIdleConns = value
	return t
}

// MaxIdleConnsPerHost sets http.Transport.MaxIdleConnsPerHost.
func (t *Transport) MaxIdleConnsPerHost(value int) *Transport {
	t.http.MaxConnsPerHost = value
	return t
}

// MaxConnsPerHost sets http.Transport.MaxConnsPerHost.
func (t *Transport) MaxConnsPerHost(value int) *Transport {
	t.http.MaxConnsPerHost = value
	return t
}

// IdleConnTimeout sets http.Transport.IdleConnTimeout.
func (t *Transport) IdleConnTimeout(value time.Duration) *Transport {
	t.http.IdleConnTimeout = value
	return t
}

// ResponseHeaderTimeout sets http.Transport.ResponseHeaderTimeout.
func (t *Transport) ResponseHeaderTimeout(value time.Duration) *Transport {
	t.http.ResponseHeaderTimeout = value
	return t
}

// TLSHandshakeTimeout sets http.Transport.TLSHandshakeTimeout.
func (t *Transport) TLSHandshakeTimeout(value time.Duration) *Transport {
	t.http.TLSHandshakeTimeout = value
	return t
}

func (t *Transport) TLSClientConfig(value *tls.Config) *Transport {
	t.http.TLSClientConfig = value
	return t
}

// ExpectContinueTimeout sets http.Transport.ExpectContinueTimeout.
func (t *Transport) ExpectContinueTimeout(value time.Duration) *Transport {
	t.http.ExpectContinueTimeout = value
	return t
}

func (t *Transport) Restraint(limiter Limiter) *Transport {
	t.limiter = limiter
	return t
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.limiter.Start()
	defer t.limiter.Complete()
	return t.http.RoundTrip(req)
}

// NewClient creates a new Client with this Transport.
func (t *Transport) NewClient() *Client {
	return NewClient(&http.Client{Transport: t})
}
