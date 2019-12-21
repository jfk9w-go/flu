package flu

import "io"

type Readable interface {
	Reader() (io.ReadCloser, error)
}

type Writable interface {
	Writer() (io.WriteCloser, error)
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

func Write(writer WriterTo, writable Writable) error {
	w, err := writable.Writer()
	if err != nil {
		return err
	}
	//noinspection GoUnhandledErrorResult
	defer w.Close()
	return writer.WriteTo(w)
}

func Read(readable Readable, reader ReaderFrom) error {
	r, err := readable.Reader()
	if err != nil {
		return err
	}
	//noinspection GoUnhandledErrorResult
	defer r.Close()
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
func Copy(readable Readable, writable Writable) error {
	reader, err := readable.Reader()
	if err != nil {
		return err
	}
	defer reader.Close()
	writer, err := writable.Writer()
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = io.Copy(writer, reader)
	return err
}
