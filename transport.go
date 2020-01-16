package flu

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Transport is a fluent wrapper around *http.Transport.
type Transport struct {
	http      *http.Transport
	logger    *log.Logger
	requests  *sync.Map
	restraint Restraint
}

// NewTransport initializes a new Transport with default settings.
// This should be equivalent to http.DefaultTransport
func NewTransport() *Transport {
	return &Transport{
		http:      http.DefaultTransport.(*http.Transport).Clone(),
		logger:    nil,
		requests:  nil,
		restraint: NoRestraint,
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

func (t *Transport) Logger(logger *log.Logger) *Transport {
	t.logger = logger
	return t
}

func (t *Transport) PendingRequests(requests *sync.Map) *Transport {
	t.requests = requests
	return t
}

func (t *Transport) Restraint(restraint Restraint) *Transport {
	t.restraint = restraint
	return t
}

var RequestLogIDLength = 8

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var id string
	var description string
	var startTime time.Time
	if t.logger != nil || t.requests != nil {
		id = GenerateID(RequestLogIDLength)
		description = req.Method + " " + req.URL.String()
		if t.logger != nil {
			startTime = time.Now()
			t.logger.Printf("[%s] %s ...", id, description)
		}
		if t.requests != nil {
			t.requests.Store(id, description)
		}
	}
	t.restraint.Start()
	resp, err := t.http.RoundTrip(req)
	t.restraint.Complete()
	if t.logger != nil {
		duration := time.Now().Sub(startTime)
		if err != nil {
			t.logger.Printf("[%s] %s %s %s (%v)", id, req.Method, req.URL.String(), err, duration)
		} else {
			t.logger.Printf("[%s] %s %s %s (%v)",
				id, req.Method, req.URL.String(), resp.Status, duration)
		}
	}
	if t.requests != nil {
		t.requests.Delete(id)
	}
	return resp, err
}

// NewClient creates a new Client with this Transport.
func (t *Transport) NewClient() *Client {
	return NewClient(&http.Client{Transport: t})
}
