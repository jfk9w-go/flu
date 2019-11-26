package flu

import (
	"io"
	"net/url"

	"github.com/google/go-querystring/query"
)

type Form struct {
	v  interface{}
	uv url.Values
}

func FormValue(v interface{}) Form {
	return Form{v: v}
}

func FormValues(uv url.Values) Form {
	return Form{uv: uv}
}

func EmptyForm() Form {
	return Form{}
}

func (f Form) EncodeTo(w io.Writer) error {
	uv, err := f.encodeValue()
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, uv.Encode())
	if err != nil {
		return err
	}

	if len(f.uv) > 0 && len(uv) > 0 {
		_, err = io.WriteString(w, "&")
		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(w, f.uv.Encode())
	return err
}

func (f Form) ContentType() string {
	return "application/x-www-form-urlencoded"
}

func (f Form) Set(k, v string) Form {
	f.uv.Set(k, v)
	return f
}

func (f Form) Add(k, v string) Form {
	f.uv.Add(k, v)
	return f
}

func (f Form) AddAll(k string, vs ...string) Form {
	for _, v := range vs {
		f.Add(k, v)
	}

	return f
}

func (f Form) Multipart() MultipartForm {
	return MultipartFormFrom(f)
}

func (f Form) encodeValue() (url.Values, error) {
	uv, err := query.Values(f.v)
	if err != nil {
		return nil, err
	}

	for k := range f.uv {
		uv.Del(k)
	}

	return uv, nil
}
