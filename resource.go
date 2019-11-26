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

type FileResource string

func File(path string) FileResource {
	return FileResource(path)
}

func (r FileResource) Path() string {
	return string(r)
}

func (r FileResource) Reader() (io.ReadCloser, error) {
	return os.Open(r.Path())
}

func (r FileResource) Writer() (io.WriteCloser, error) {
	if err := os.MkdirAll(path.Dir(r.Path()), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(r.Path())
}

type URLResource string

func URL(rawurl string) URLResource {
	return URLResource(rawurl)
}

func (r URLResource) URL() string {
	return string(r)
}

func (r URLResource) Reader() (io.ReadCloser, error) {
	resp, err := http.Get(r.URL())
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
