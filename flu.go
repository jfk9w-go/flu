package flu

import (
	"io"
)

type Readable interface {
	Reader() (io.Reader, error)
}

type Writable interface {
	Writer() (io.Writer, error)
}

type ReadWritable interface {
	Readable
	Writable
}

type EncoderTo interface {
	EncodeTo(io.Writer) error
}

type DecoderFrom interface {
	DecodeFrom(io.Reader) error
}

type Body interface {
	ContentType() string
}

type BodyEncoderTo interface {
	Body
	EncoderTo
}

type BodyDecoderFrom interface {
	Body
	DecoderFrom
}

type BodyCodec interface {
	Body
	EncoderTo
	DecoderFrom
}

func EncodeTo(encoder EncoderTo, out Writable) error {
	writer, err := out.Writer()
	if err != nil {
		return err
	}
	if closer, ok := writer.(io.Closer); ok {
		defer closer.Close()
	}
	return encoder.EncodeTo(writer)
}

func DecodeFrom(in Readable, decoder DecoderFrom) error {
	reader, err := in.Reader()
	if err != nil {
		return err
	}
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}
	return decoder.DecodeFrom(reader)
}

func AsReadable(encoder EncoderTo) Readable {
	reader, writer := io.Pipe()
	go func() {
		err := encoder.EncodeTo(writer)
		_ = writer.CloseWithError(err)
	}()
	return Xable{R: reader}
}

func AsWritable(decoder DecoderFrom) Writable {
	reader, writer := io.Pipe()
	go func() {
		err := decoder.DecodeFrom(reader)
		_ = reader.CloseWithError(err)
	}()
	return Xable{W: writer}
}

func Copy(in Readable, out Writable) error {
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
