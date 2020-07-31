package flu

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"path"
)

type IO struct {
	R io.Reader
	W io.Writer
}

func (io IO) Reader() (io.Reader, error) {
	return io.R, nil
}

func (io IO) Writer() (io.Writer, error) {
	return io.W, nil
}

type File string

func (f File) Path() string {
	return string(f)
}

func (f File) Open() (*os.File, error) {
	return os.Open(f.Path())
}

func (f File) Create() (*os.File, error) {
	if err := os.MkdirAll(path.Dir(f.Path()), os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(f.Path())
}

func (f File) Reader() (io.Reader, error) {
	return f.Open()
}

func (f File) Writer() (io.Writer, error) {
	return f.Create()
}

type URL string

func (u URL) URL() string {
	return string(u)
}

func (u URL) Reader() (io.Reader, error) {
	resp, err := http.Get(string(u))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() Buffer {
	return Buffer{new(bytes.Buffer)}
}

func (b Buffer) Reader() (io.Reader, error) {
	return Bytes(b.Bytes()).Reader()
}

func (b Buffer) Writer() (io.Writer, error) {
	b.Reset()
	return b.Buffer, nil
}

type Bytes []byte

func (b Bytes) Reader() (io.Reader, error) {
	return bytes.NewReader(b), nil
}

type Conn struct {
	Dialer  net.Dialer
	Context context.Context
	Network string
	Address string
}

func (c Conn) Dial() (net.Conn, error) {
	if c.Context != nil {
		return c.Dialer.DialContext(c.Context, c.Network, c.Address)
	} else {
		return c.Dialer.Dial(c.Network, c.Address)
	}
}

func (c Conn) Reader() (io.Reader, error) {
	return c.Dial()
}

func (c Conn) Writer() (io.Writer, error) {
	return c.Dial()
}
