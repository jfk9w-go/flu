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
	writer, err := out.Writer()
	if err != nil {
		return err
	}
	if closer, ok := writer.(io.Closer); ok {
		defer closer.Close()
	}
	return encoder.EncodeTo(writer)
}

func DecodeFrom(in Input, decoder DecoderFrom) error {
	reader, err := in.Reader()
	if err != nil {
		return err
	}
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}
	return decoder.DecodeFrom(reader)
}

func ReadablePipe(encoder EncoderTo) Input {
	reader, writer := io.Pipe()
	go func() {
		err := encoder.EncodeTo(writer)
		_ = writer.CloseWithError(err)
	}()
	return IO{R: reader}
}

func WritablePipe(decoder DecoderFrom) Output {
	reader, writer := io.Pipe()
	go func() {
		err := decoder.DecodeFrom(reader)
		_ = reader.CloseWithError(err)
	}()
	return IO{W: writer}
}

func Copy(in Input, out Output) error {
	reader, err := in.Reader()
	if err != nil {
		return err
	}
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}
	writer, err := out.Writer()
	if err != nil {
		return err
	}
	if closer, ok := writer.(io.Closer); ok {
		defer closer.Close()
	}
	_, err = io.Copy(writer, reader)
	return err
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
