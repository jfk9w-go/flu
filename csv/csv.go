package csv

import (
	"encoding/csv"
	"io"

	"golang.org/x/text/encoding"

	"github.com/pkg/errors"
)

type ParseError struct {
	Values []string
	Err    error
}

func (e ParseError) Error() string {
	return e.Err.Error()
}

type Codec struct {
	Encoding encoding.Encoding
	Output   RowOutput
	Header   []string
	Comma    rune
}

func (c *Codec) DecodeFrom(r io.Reader) error {
	if c.Encoding != nil {
		r = c.Encoding.NewDecoder().Reader(r)
	}

	csv := csv.NewReader(r)
	csv.Comma = c.Comma
	if c.Header == nil {
		var err error
		if c.Header, err = csv.Read(); err != nil {
			return errors.Wrap(err, "parse header")
		}
	}

	for {
		values, err := csv.Read()
		if err == nil && len(c.Header) != len(values) {
			err = errors.Errorf(
				"header size (%d) does not match row size (%d)",
				len(c.Header), len(values))
		}

		switch err {
		case io.EOF:
			return nil
		case nil:
			row := &Row{Values: make(map[string]string, len(c.Header))}
			for i, key := range c.Header {
				row.Values[key] = values[i]
			}

			err = c.Output.Output(row)
		default:
			err = c.Output.Output(&Row{Err: ParseError{
				Values: values,
				Err:    err,
			}})
		}

		if err != nil {
			return errors.Wrap(err, "output csv row")
		}
	}
}
