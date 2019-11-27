package flu

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

func FormValue(value interface{}, withValues bool) Form {
	var values url.Values = nil
	if withValues {
		values = make(url.Values)
	}

	return Form{value: value, values: values}
}

func FormValues(values url.Values) Form {
	return Form{values: values}
}

func EmptyForm(withValues bool) Form {
	return FormValue(nil, withValues)
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
	form.values.Set(key, value)
	return form
}

func (form Form) Add(key, value string) Form {
	form.values.Add(key, value)
	return form
}

func (form Form) AddAll(key string, values ...string) Form {
	for _, v := range values {
		form.Add(key, v)
	}

	return form
}

func (form Form) Multipart() MultipartForm {
	return MultipartFormFrom(form)
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
