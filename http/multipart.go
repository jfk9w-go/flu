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
	*Form
	boundary string
	files    map[string]multipartFile
}

func NewMultipartForm() *MultipartForm {
	return &MultipartForm{
		Form:     new(Form),
		boundary: randomBoundary(),
	}
}

func (mf *MultipartForm) Set(key, value string) *MultipartForm {
	mf.Form = mf.Form.Set(key, value)
	return mf
}

func (mf *MultipartForm) Add(key, value string) *MultipartForm {
	mf.Form = mf.Form.Add(key, value)
	return mf
}

func (mf *MultipartForm) AddAll(keys string, values ...string) *MultipartForm {
	mf.Form = mf.Form.AddAll(keys, values...)
	return mf
}

func (mf *MultipartForm) SetValues(values url.Values) *MultipartForm {
	mf.Form = mf.Form.SetValues(values)
	return mf
}

func (mf *MultipartForm) AddValues(values url.Values) *MultipartForm {
	mf.Form = mf.Form.AddValues(values)
	return mf
}

func (mf *MultipartForm) Value(value interface{}) *MultipartForm {
	mf.Form = mf.Form.Value(value)
	return mf
}

func (mf *MultipartForm) File(fieldname, filename string, input flu.Input) *MultipartForm {
	if mf.files == nil {
		mf.files = make(map[string]multipartFile)
	}
	if filename == "" {
		filename = fieldname
	}
	mf.files[fieldname] = multipartFile{
		name:  filename,
		input: input,
	}
	return mf
}

func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

func (mf *MultipartForm) EncodeTo(w io.Writer) error {
	mw := multipart.NewWriter(w)
	defer mw.Close()
	err := mw.SetBoundary(mf.boundary)
	if err != nil {
		return err
	}
	for fieldname, file := range mf.files {
		w, err := mw.CreateFormFile(fieldname, file.name)
		if err != nil {
			return err
		}
		_, err = flu.Copy(file.input, flu.IO{W: w})
		if err != nil {
			return err
		}
	}
	values, err := mf.Form.encodeValue()
	if err != nil {
		return err
	}
	err = writeMultipartValues(mw, values)
	if err != nil {
		return err
	}
	return writeMultipartValues(mw, mf.values)
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

func (mf *MultipartForm) ContentType() string {
	return "multipart/form-data; boundary=" + mf.boundary
}

type multipartFile struct {
	name  string
	input flu.Input
}
