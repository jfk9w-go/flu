package flu

import (
	"crypto/rand"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

type MultipartForm struct {
	Form
	b  string
	rs map[string]ResourceReader
}

func EmptyMultipartForm(withFormValues bool) MultipartForm {
	return MultipartFormFrom(EmptyForm(withFormValues))
}

func MultipartFormFrom(f Form) MultipartForm {
	return MultipartForm{
		Form: f,
		b:    randomBoundary(),
		rs:   make(map[string]ResourceReader),
	}
}

func MultipartFormValues(uv url.Values) MultipartForm {
	return MultipartFormFrom(FormValues(uv))
}

func (f MultipartForm) Set(k, v string) MultipartForm {
	f.Form.Set(k, v)
	return f
}

func (f MultipartForm) Add(k, v string) MultipartForm {
	f.Form.Add(k, v)
	return f
}

func (f MultipartForm) AddAll(k string, vs ...string) MultipartForm {
	f.Form.AddAll(k, vs...)
	return f
}

func (f MultipartForm) Resource(k string, r ResourceReader) MultipartForm {
	f.rs[k] = r
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
	multipartWriter := multipart.NewWriter(w)
	//noinspection GoUnhandledErrorResult
	defer multipartWriter.Close()

	err := multipartWriter.SetBoundary(f.b)
	if err != nil {
		return err
	}

	for key, resource := range f.rs {
		w, err := multipartWriter.CreateFormFile(key, key)
		if err != nil {
			return err
		}

		r, err := resource.Reader()
		if err != nil {
			return err
		}

		_, err = io.Copy(w, r)
		_ = r.Close()
		if err != nil {
			return err
		}
	}

	urlValues, err := f.Form.encodeValue()
	if err != nil {
		return err
	}

	err = writeMultipartValues(multipartWriter, urlValues)
	if err != nil {
		return err
	}

	return writeMultipartValues(multipartWriter, f.values)
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
	return "multipart/form-data; boundary=" + f.b
}
