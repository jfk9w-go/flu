package http

import (
	"crypto/rand"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"

	"github.com/jfk9w-go/flu"
)

type MultipartForm struct {
	Form
	boundary string
	files    map[string]flu.Input
}

func NewMultipartForm() MultipartForm {
	return MultipartForm{
		boundary: randomBoundary(),
	}
}

func (f MultipartForm) Set(k, v string) MultipartForm {
	f.Form = f.Form.Set(k, v)
	return f
}

func (f MultipartForm) Add(k, v string) MultipartForm {
	f.Form = f.Form.Add(k, v)
	return f
}

func (f MultipartForm) AddAll(k string, vs ...string) MultipartForm {
	f.Form = f.Form.AddAll(k, vs...)
	return f
}

func (f MultipartForm) SetValues(values url.Values) MultipartForm {
	f.Form = f.Form.SetValues(values)
	return f
}

func (f MultipartForm) AddValues(values url.Values) MultipartForm {
	f.Form = f.Form.AddValues(values)
	return f
}

func (f MultipartForm) Value(value interface{}) MultipartForm {
	f.Form = f.Form.Value(value)
	return f
}

func (f MultipartForm) File(k string, r flu.Input) MultipartForm {
	if f.files == nil {
		f.files = make(map[string]flu.Input)
	}

	f.files[k] = r
	return f
}

func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

func (f MultipartForm) EncodeTo(w io.Writer) error {
	mw := multipart.NewWriter(w)
	defer mw.Close()
	err := mw.SetBoundary(f.boundary)
	if err != nil {
		return err
	}
	for k, r := range f.files {
		w, err := mw.CreateFormFile(k, k)
		if err != nil {
			return err
		}
		err = flu.Copy(r, flu.IO{W: w})
		if err != nil {
			return err
		}
	}
	values, err := f.Form.encodeValue()
	if err != nil {
		return err
	}
	err = writeMultipartValues(mw, values)
	if err != nil {
		return err
	}
	return writeMultipartValues(mw, f.values)
}

func writeMultipartValues(mw *multipart.Writer, uv url.Values) error {
	for k, vs := range uv {
		for _, value := range vs {
			err := mw.WriteField(k, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (f MultipartForm) ContentType() string {
	return "multipart/form-data; boundary=" + f.boundary
}
