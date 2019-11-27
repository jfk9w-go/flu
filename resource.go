package flu

import (
	"io"
	"net/http"
	"os"
	"path"
)

type ResourceReader interface {
	Reader() (io.ReadCloser, error)
}

type ResourceWriter interface {
	Writer() (io.WriteCloser, error)
}

type File string

func (f File) Path() string {
	return string(f)
}

func (f File) Reader() (io.ReadCloser, error) {
	return os.Open(f.Path())
}

func (f File) Writer() (io.WriteCloser, error) {
	if err := os.MkdirAll(path.Dir(f.Path()), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(f.Path())
}

type URL string

func (u URL) URL() string {
	return string(u)
}

func (u URL) Reader() (io.ReadCloser, error) {
	resp, err := http.Get(u.URL())
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
