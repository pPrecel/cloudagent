package output

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"github.com/lensesio/tableprinter"
	"github.com/tidwall/gjson"
)

type TextOptions struct {
	Format string
	RPath  string
	HPath  string
	UPath  string
	APath  string
}

func PrintText(w io.Writer, v interface{}, o TextOptions, f ...string) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	s := o.Format
	s = strings.Replace(s, "%a", gjson.GetBytes(b, o.APath).Raw, -1)

	for i := range f {
		b = []byte(gjson.GetBytes(b, f[i]).Raw)
	}

	s = strings.Replace(s, "%u", gjson.GetBytes(b, o.UPath).Raw, -1)
	s = strings.Replace(s, "%h", gjson.GetBytes(b, o.HPath).Raw, -1)
	s = strings.Replace(s, "%r", gjson.GetBytes(b, o.RPath).Raw, -1)

	_, err = w.Write([]byte(s))

	return err
}

type ErrorOptions struct {
	Format string
	Error  string
}

func PrintErrorText(w io.Writer, o ErrorOptions) error {
	s := o.Format
	s = strings.Replace(s, "%e", o.Error, -1)

	_, err := w.Write([]byte(s))

	return err
}

func PrintJson(w io.Writer, v interface{}, f ...string) error {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}

	for i := range f {
		b = []byte(gjson.GetBytes(b, f[i]).Raw)
	}

	_, err = w.Write(b)

	return err
}

func PrintTable(w io.Writer, v interface{}, f ...string) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	for i := range f {
		b = []byte(gjson.GetBytes(b, f[i]).Raw)
	}

	val := reflect.New(reflect.TypeOf(v)).Interface()

	err = json.Unmarshal(b, val)
	if err != nil {
		return err
	}

	p := tableprinter.New(w)

	p.BorderTop = true
	p.BorderBottom = true
	p.BorderLeft = true
	p.BorderRight = true
	p.ColumnSeparator = "│"
	p.RowSeparator = "─"

	p.Print(val)

	return nil
}
