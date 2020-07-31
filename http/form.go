package http

import (
	"io"
	"net/url"

	"github.com/google/go-querystring/query"
)

type Form struct {
	// value contains a value which will be transformed into url.Values
	value interface{}
	// values contains values which can be set separately by key
	// keys from here override the same keys in url.Values transformed from values
	values url.Values
}

func (f *Form) EncodeTo(writer io.Writer) error {
	values, err := f.encodeValue()
	if err != nil {
		return err
	}
	_, err = io.WriteString(writer, values.Encode())
	if err != nil {
		return err
	}
	if len(f.values) > 0 && len(values) > 0 {
		_, err = io.WriteString(writer, "&")
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(writer, f.values.Encode())
	return err
}

func (*Form) ContentType() string {
	return "application/x-www-form-urlencoded"
}

func (f *Form) Set(key, value string) *Form {
	if f.values == nil {
		f.values = make(url.Values)
	}

	f.values.Set(key, value)
	return f
}

func (f *Form) Add(key, value string) *Form {
	if f.values == nil {
		f.values = make(url.Values)
	}

	f.values.Add(key, value)
	return f
}

func (f *Form) AddAll(key string, values ...string) *Form {
	for _, v := range values {
		f.Add(key, v)
	}

	return f
}

func (f *Form) SetValues(values url.Values) *Form {
	if f.values == nil {
		f.values = values
	} else {
		a, b := f.values, values
		if len(b) < len(a) {
			a, b = b, a
		}
		for k, v := range a {
			b[k] = v
		}
	}

	return f
}

func (f *Form) AddValues(values url.Values) *Form {
	if f.values == nil {
		f.values = values
	} else {
		a, b := f.values, values
		if len(b) < len(a) {
			a, b = b, a
		}
		for k, v := range a {
			b[k] = append(b[k], v...)
		}
	}

	return f
}

func (f *Form) Value(value interface{}) *Form {
	f.value = value
	return f
}

func (f *Form) Multipart() *MultipartForm {
	multipart := NewMultipartForm()
	multipart.Form = f
	return multipart
}

func (f *Form) encodeValue() (url.Values, error) {
	values, err := query.Values(f.value)
	if err != nil {
		return nil, err
	}
	for k := range f.values {
		values.Del(k)
	}
	return values, nil
}
