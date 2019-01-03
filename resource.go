package flu

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

// ReadResource provides an access to a resource which can be read.
type ReadResource interface {
	Read() (io.ReadCloser, error)
}

// WriteResource provides an access to a resource which can be written.
type WriteResource interface {
	Write() (io.WriteCloser, error)
}

// ReadWriteResource provides full read-write access to a resource.
type ReadWriteResource interface {
	ReadResource
	WriteResource
}

// RawReadResource is a wrapper around io.ReadCloser.
// Provides "read-once" semantics.
type RawReadResource struct {
	rc io.ReadCloser
}

// NewReadResource wrappes a io.Reader into a Resource.
func NewReadResource(reader io.Reader) RawReadResource {
	var (
		rc io.ReadCloser
		ok = false
	)

	if rc, ok = reader.(io.ReadCloser); !ok {
		rc = ioutil.NopCloser(reader)
	}

	return RawReadResource{rc}
}

// Read returns the wrapped io.Reader.
func (r RawReadResource) Read() (io.ReadCloser, error) {
	return r.rc, nil
}

// FileSystemResource is a file identified by its path.
type FileSystemResource struct {
	path string
}

// NewFileSystemResource creates a FileSystemResource with the specified path.
func NewFileSystemResource(path string) *FileSystemResource {
	return &FileSystemResource{path}
}

// Read opens the file for reading.
func (r *FileSystemResource) Read() (io.ReadCloser, error) {
	return os.Open(r.path)
}

// Write creates (with all the intermediary folders) or truncates the file
// and opens it for writing.
func (r *FileSystemResource) Write() (io.WriteCloser, error) {
	if err := os.MkdirAll(path.Dir(r.path), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(r.path)
}

// Size returns the file size.
func (r *FileSystemResource) Size() (size int64, err error) {
	stat, err := os.Stat(r.path)
	if err != nil {
		return
	}

	size = stat.Size()
	return
}

// Delete deletes the file.
func (r *FileSystemResource) Delete() error {
	return os.RemoveAll(r.path)
}

// Path returns the file path.
func (r *FileSystemResource) Path() string {
	return r.path
}
