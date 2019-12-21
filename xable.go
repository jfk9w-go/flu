package flu

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type Xable struct {
	R io.Reader
	W io.Writer
}

func (x Xable) Reader() (io.ReadCloser, error) {
	if rc, ok := x.R.(io.ReadCloser); ok {
		return rc, nil
	} else {
		return ioutil.NopCloser(x.R), nil
	}
}

func (x Xable) Writer() (io.WriteCloser, error) {
	return x, nil
}

func (x Xable) Write(p []byte) (int, error) {
	return x.W.Write(p)
}

func (x Xable) Close() error {
	if wc, ok := x.W.(io.Closer); ok {
		return wc.Close()
	} else {
		return nil
	}
}

type File string

func (f File) Path() string {
	return string(f)
}

func (f File) Reader() (io.ReadCloser, error) {
	return os.Open(f.Path())
}

func (f File) Writer() (io.WriteCloser, error) {
	if err := os.MkdirAll(path.Dir(f.Path()), os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(f.Path())
}

type URL string

func (u URL) URL() string {
	return string(u)
}

func (u URL) Reader() (io.ReadCloser, error) {
	resp, err := http.Get(u.URL())
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

type Buffer bytes.Buffer

func (b *Buffer) bb() *bytes.Buffer {
	return (*bytes.Buffer)(b)
}

func (b *Buffer) Reader() (io.ReadCloser, error) {
	return Bytes(b.bb().Bytes()).Reader()
}

func (b *Buffer) ReadFrom(r io.Reader) error {
	return Copy(Xable{R: r}, b)
}

func (b *Buffer) Writer() (io.WriteCloser, error) {
	b.bb().Reset()
	return b, nil
}

func (b *Buffer) Write(p []byte) (int, error) {
	return b.bb().Write(p)
}

func (b *Buffer) Close() error {
	return nil
}

type Bytes []byte

func (b Bytes) Reader() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader(b)), nil
}
