package csv

import (
	"strconv"
	"strings"
	"time"
)

type Value interface {
	Parse(str string) error
}

type String string

func (s *String) Parse(str string) error {
	*(*string)(s) = str
	return nil
}

func (s String) Value() string {
	return string(s)
}

var DefaultFloatFormat = &FloatFormat{
	Replacer: strings.NewReplacer(" ", "", ",", "."),
}

type FloatFormat struct {
	Replacer *strings.Replacer
}

type Float struct {
	Value   float64
	BitSize int
	Format  *FloatFormat
}

func (f *Float) Parse(str string) error {
	format := f.Format
	if format == nil {
		format = DefaultFloatFormat
	}

	str = format.Replacer.Replace(str)
	var err error
	f.Value, err = strconv.ParseFloat(str, f.BitSize)
	return err
}

type Time struct {
	Value    time.Time
	Layout   string
	Location *time.Location
}

func (t *Time) Parse(str string) error {
	var err error
	t.Value, err = time.ParseInLocation(t.Layout, str, t.Location)
	return err
}

type Uint struct {
	Value         uint64
	Base, BitSize int
}

func (u *Uint) Parse(str string) error {
	var err error
	u.Value, err = strconv.ParseUint(str, u.Base, u.BitSize)
	return err
}
