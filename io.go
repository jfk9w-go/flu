package flu

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"
)

type Input interface {
	Reader() (io.Reader, error)
}

type Output interface {
	Writer() (io.Writer, error)
}

type InputOutput interface {
	Input
	Output
}

type EncoderTo interface {
	EncodeTo(io.Writer) error
}

type DecoderFrom interface {
	DecodeFrom(io.Reader) error
}

func EncodeTo(encoder EncoderTo, out Output) error {
	w, err := out.Writer()
	if err != nil {
		return err
	}
	if c, ok := w.(io.Closer); ok {
		defer c.Close()
	}
	return encoder.EncodeTo(w)
}

func DecodeFrom(in Input, decoder DecoderFrom) error {
	r, err := in.Reader()
	if err != nil {
		return err
	}
	if c, ok := r.(io.Closer); ok {
		defer c.Close()
	}
	return decoder.DecodeFrom(r)
}

func PipeInput(encoder EncoderTo) Input {
	r, w := io.Pipe()
	go func() {
		err := encoder.EncodeTo(w)
		_ = w.CloseWithError(err)
	}()
	return IO{R: r}
}

func PipeOutput(decoder DecoderFrom) Output {
	r, w := io.Pipe()
	go func() {
		err := decoder.DecodeFrom(r)
		_ = r.CloseWithError(err)
	}()
	return IO{W: w}
}

func Copy(in Input, out Output) error {
	r, err := in.Reader()
	if err != nil {
		return err
	}
	if c, ok := r.(io.Closer); ok {
		defer c.Close()
	}
	w, err := out.Writer()
	if err != nil {
		return err
	}
	if c, ok := w.(io.Closer); ok {
		defer c.Close()
	}
	_, err = io.Copy(w, r)
	return err
}

type IOCounter int64

func (c *IOCounter) Write(data []byte) (n int, err error) {
	n = len(data)
	*(*int64)(c) += int64(n)
	return n, nil
}

func (c *IOCounter) Count(encoder EncoderTo) error {
	return EncodeTo(encoder, IO{W: c})
}

func (c *IOCounter) Value() int64 {
	return *(*int64)(c)
}

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
