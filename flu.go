package flu

import "io"

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

type WriterTo interface {
	WriteTo(io.Writer) error
}

type ReaderFrom interface {
	ReadFrom(io.Reader) error
}

type Body interface {
	ContentType() string
}

type BodyWriter interface {
	Body
	WriterTo
}

type BodyReader interface {
	Body
	ReaderFrom
}

type BodyReadWriter interface {
	Body
	WriterTo
	ReaderFrom
}

func Write(writer WriterTo, out Writable) error {
	w, err := out.Writer()
	if err != nil {
		return err
	}
	//noinspection GoUnhandledErrorResult
	if w, ok := w.(io.Closer); ok {
		defer w.Close()
	}
	return writer.WriteTo(w)
}

func Read(in Readable, reader ReaderFrom) error {
	r, err := in.Reader()
	if err != nil {
		return err
	}
	//noinspection GoUnhandledErrorResult
	if r, ok := r.(io.Closer); ok {
		defer r.Close()
	}
	return reader.ReadFrom(r)
}

func PipeOut(writer WriterTo) Readable {
	out, in := io.Pipe()
	go func() {
		err := writer.WriteTo(in)
		_ = in.CloseWithError(err)
	}()
	return Xable{R: out}
}

func PipeIn(reader ReaderFrom) Writable {
	out, in := io.Pipe()
	go func() {
		err := reader.ReadFrom(out)
		_ = out.CloseWithError(err)
	}()
	return Xable{W: in}
}

//noinspection GoUnhandledErrorResult
func Copy(in Readable, out Writable) error {
	r, err := in.Reader()
	if err != nil {
		return err
	}
	if r, ok := r.(io.Closer); ok {
		defer r.Close()
	}
	w, err := out.Writer()
	if err != nil {
		return err
	}
	if w, ok := w.(io.Closer); ok {
		defer w.Close()
	}
	_, err = io.Copy(w, r)
	return err
}

var DefaultClient = NewClient(nil)

func NewRequest(resource string) *Request {
	return DefaultClient.NewRequest(resource)
}
