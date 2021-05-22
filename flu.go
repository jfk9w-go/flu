package flu

import (
	"context"
	"io"
	"log"
	"time"
)

// Input interface describes a resource which can be read
// (possibly more than once).
type Input interface {
	// Reader returns an instance of io.Reader.
	Reader() (io.Reader, error)
}

// Writer interface describes a resource which can be written
// (possibly more than once).
type Output interface {
	// Writer returns an instance of io.Writer.
	Writer() (io.Writer, error)
}

// EncoderTo interface describes a value which can be encoded.
type EncoderTo interface {
	// EncodeTo encodes the value to the given io.Writer.
	EncodeTo(io.Writer) error
}

// DecoderFrom interface describes a value which can be decoded.
type DecoderFrom interface {
	// DecodeFrom decodes the value from the given io.Reader.
	DecodeFrom(io.Reader) error
}

// EncodeTo encodes the provided EncoderTo to Output.
// It closes the io.Writer instance if necessary.
func EncodeTo(encoder EncoderTo, out Output) error {
	w, err := out.Writer()
	if err != nil {
		return err
	}
	if err := encoder.EncodeTo(w); err != nil {
		return err
	}
	return Close(w)
}

// DecodeFrom decodes the provided DecoderFrom from Input.
// It closes the io.Reader instance if necessary.
func DecodeFrom(in Input, decoder DecoderFrom) error {
	r, err := in.Reader()
	if err != nil {
		return err
	}
	if err := decoder.DecodeFrom(r); err != nil {
		return err
	}
	return Close(r)
}

// PipeInput pipes the encoded value from EncoderTo as Input
// in the background.
func PipeInput(encoder EncoderTo) Input {
	r, w := io.Pipe()
	go func() {
		err := encoder.EncodeTo(w)
		if err := w.CloseWithError(err); err != nil {
			log.Printf("PipeInput close error: %s", err)
		}
	}()

	return IO{R: r}
}

// PipeOutput provides an Output which feeds into DecoderFrom
// in the background.
func PipeOutput(decoder DecoderFrom) Output {
	r, w := io.Pipe()
	go func() {
		err := decoder.DecodeFrom(r)
		if err := r.CloseWithError(err); err != nil {
			log.Printf("PipeOutput close error: %s", err)
		}
	}()

	return IO{W: w}
}

// Copy copies the Input to the Output.
func Copy(in Input, out Output) (written int64, err error) {
	r, err := in.Reader()
	if err != nil {
		return
	}
	w, err := out.Writer()
	if err != nil {
		return
	}
	written, err = io.Copy(w, r)
	if err != nil {
		return
	}
	if err = Close(w); err != nil {
		return
	}
	if err = Close(r); err != nil {
		return
	}
	return
}

// Sleep sleeps for the specified timeout interruptibly.
func Sleep(ctx context.Context, timeout time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(timeout):
		return nil
	}
}
