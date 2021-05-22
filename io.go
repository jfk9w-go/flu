package flu

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/text/encoding"
)

// IO is a generic Input / Output.
// It is not mandatory to fill all struct fields.
type IO struct {
	// R is an io.Reader instance to be used for reading.
	R io.Reader
	// W is an io.Writer instance to be used for reading.
	W io.Writer
	// E is an error which will be returned when reading/writing.
	E error
}

func (io IO) Reader() (io.Reader, error) {
	if io.E != nil {
		return nil, io.E
	}

	return io.R, nil
}

func (io IO) Writer() (io.Writer, error) {
	if io.E != nil {
		return nil, io.E
	}

	return io.W, nil
}

// File is a path representing a file (or directory).
type File string

// FilePath creates a File instance from the provided
// path parts.
func FilePath(path ...string) File {
	return File(filepath.Join(path...))
}

// Path returns the underlying string.
func (f File) Path() string {
	return string(f)
}

// Join creates a new File instance pointing
// to the child element of this instance.
func (f File) Join(child string) File {
	return FilePath(f.Path(), child)
}

// Exists checks for the existence of the File entry.
func (f File) Exists() (bool, error) {
	_, err := os.Stat(f.Path())
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

// Open opens the File for reading.
func (f File) Open() (*os.File, error) {
	return os.Open(f.Path())
}

// Create opens the File for writing.
// It creates the file and all intermediate directories if necessary.
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

// Remove removes the file or directory represented by this File.
func (f File) Remove() error {
	return os.RemoveAll(f.Path())
}

// URL is a read-only resource accessible by URL.
type URL string

// Unmask returns the underlying string.
func (u URL) Unmask() string {
	return string(u)
}

func (u URL) Reader() (io.Reader, error) {
	resp, err := http.Get(string(u))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// ByteBuffer is an Input / Output wrapper around bytes.Buffer.
type ByteBuffer bytes.Buffer

func (b *ByteBuffer) Reader() (io.Reader, error) {
	return b.Unmask(), nil
}

func (b *ByteBuffer) Writer() (io.Writer, error) {
	b.Unmask().Reset()
	return b.Unmask(), nil
}

// Bytes returns read-only Bytes view on this buffer.
func (b *ByteBuffer) Bytes() Bytes {
	return Bytes(b.Unmask().Bytes())
}

// Unmask returns the underlying *bytes.Buffer.
func (b *ByteBuffer) Unmask() *bytes.Buffer {
	return (*bytes.Buffer)(b)
}

// Bytes is a read-only byte array.
type Bytes []byte

func (b Bytes) Reader() (io.Reader, error) {
	return bytes.NewReader(b), nil
}

// Conn provides the means for opening net.Conn.
type Conn struct {
	// Dialer is the net.Dialer to be used for connection.
	// May be empty.
	Dialer net.Dialer

	// Context is the context.Context to be used.
	// May be empty.
	Context context.Context

	// Network is the network passed to Dialer.Dial.
	Network string

	// Address is the address passed to Dialer.Dial.
	Address string
}

// Dial opens a net.Conn using the provided struct fields.
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

// Close attempts to close the provided value
// using io.Closer interface.
func Close(value interface{}) error {
	if closer, ok := value.(io.Closer); ok {
		return closer.Close()
	} else {
		return nil
	}
}

// AnyCloser wraps the provided value with io.Closer interface.
type AnyCloser struct {
	V interface{}
}

func (c AnyCloser) Close() error {
	return Close(c.V)
}

// Chars is the text character Input / Output wrapper.
type Chars struct {
	// In is the underlying Input.
	In Input
	// Out is the underlying Output.
	Out Output
	// Enc will be used for decoding characters from Input
	// and/or encoding them to Output.
	Enc encoding.Encoding
}

func (cs Chars) Reader() (io.Reader, error) {
	r, err := cs.In.Reader()
	if err != nil {
		return nil, err
	}

	return cs.Enc.NewDecoder().Reader(r), nil
}

func (cs Chars) Writer() (io.Writer, error) {
	w, err := cs.Out.Writer()
	if err != nil {
		return nil, err
	}

	return cs.Enc.NewEncoder().Writer(w), nil
}
