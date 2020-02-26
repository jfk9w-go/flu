package httpx

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

func (form Form) EncodeTo(writer io.Writer) error {
	values, err := form.encodeValue()
	if err != nil {
		return err
	}
	_, err = io.WriteString(writer, values.Encode())
	if err != nil {
		return err
	}
	if len(form.values) > 0 && len(values) > 0 {
		_, err = io.WriteString(writer, "&")
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(writer, form.values.Encode())
	return err
}

func (Form) ContentType() string {
	return "application/x-www-form-urlencoded"
}

func (form Form) Set(key, value string) Form {
	if form.values == nil {
		form.values = make(url.Values)
	}

	form.values.Set(key, value)
	return form
}

func (form Form) Add(key, value string) Form {
	if form.values == nil {
		form.values = make(url.Values)
	}

	form.values.Add(key, value)
	return form
}

func (form Form) AddAll(key string, values ...string) Form {
	for _, v := range values {
		form.Add(key, v)
	}

	return form
}

func (form Form) SetValues(values url.Values) Form {
	if form.values == nil {
		form.values = values
	} else {
		a, b := form.values, values
		if len(b) < len(a) {
			a, b = b, a
		}
		for k, v := range a {
			b[k] = v
		}
	}

	return form
}

func (form Form) AddValues(values url.Values) Form {
	if form.values == nil {
		form.values = values
	} else {
		a, b := form.values, values
		if len(b) < len(a) {
			a, b = b, a
		}
		for k, v := range a {
			b[k] = append(b[k], v...)
		}
	}

	return form
}

func (form Form) Value(value interface{}) Form {
	form.value = value
	return form
}

func (form Form) Multipart() MultipartForm {
	multipart := NewMultipartForm()
	multipart.Form = form
	return multipart
}

func (form Form) encodeValue() (url.Values, error) {
	values, err := query.Values(form.value)
	if err != nil {
		return nil, err
	}
	for k := range form.values {
		values.Del(k)
	}
	return values, nil
}
