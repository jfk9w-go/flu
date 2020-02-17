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
	boundary string
	files    map[string]Readable
}

func EmptyMultipartForm(withFormValues bool) MultipartForm {
	return MultipartFormFrom(EmptyForm(withFormValues))
}

func MultipartFormFrom(f Form) MultipartForm {
	return MultipartForm{
		Form:     f,
		boundary: randomBoundary(),
		files:    make(map[string]Readable),
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

func (f MultipartForm) File(k string, r Readable) MultipartForm {
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
	//noinspection GoUnhandledErrorResult
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
		err = Copy(r, Xable{W: w})
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
