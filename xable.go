package flu

import (
	"bytes"
	"io"
	"os"
	"path"
)

type Xable struct {
	R io.Reader
	W io.Writer
}

func (x Xable) Reader() (io.Reader, error) {
	return x.R, nil
}

func (x Xable) Writer() (io.Writer, error) {
	return x.W, nil
}

type File string

func (f File) Path() string {
	return string(f)
}

func (f File) Reader() (io.Reader, error) {
	return os.Open(f.Path())
}

func (f File) Writer() (io.Writer, error) {
	if err := os.MkdirAll(path.Dir(f.Path()), os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(f.Path())
}

type URL string

func (u URL) URL() string {
	return string(u)
}

func (u URL) Reader() (io.Reader, error) {
	return DefaultClient.GET(u.URL()).Execute().Reader()
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
