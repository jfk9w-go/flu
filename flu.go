package flu

import (
	"context"
	"io"
	"time"
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
	defer Close(w)
	return encoder.EncodeTo(w)
}

func DecodeFrom(in Input, decoder DecoderFrom) error {
	r, err := in.Reader()
	if err != nil {
		return err
	}
	defer Close(r)
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
	defer Close(r)
	w, err := out.Writer()
	if err != nil {
		return err
	}
	defer Close(w)
	_, err = io.Copy(w, r)
	return err
}

func Sleep(ctx context.Context, timeout time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(timeout):
		return nil
	}
}
