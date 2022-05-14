package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"
)

type OutputType string

const (
	JsonType  OutputType = "json"
	YamlType  OutputType = "yaml"
	TableType OutputType = "table"
	TextType  OutputType = "text"
)

//go:generate mockery --name=Formater --output=automock --outpkg=automock
type Formater interface {
	JSON() interface{}
	YAML() interface{}
	Table() (headers []string, data [][]string)
	Text(outFormat, errFormat string) string
}

var _ pflag.Value = &Flag{}

type Flag struct {
	outputType         OutputType
	outputStringFormat string
	outputErrorFormat  string
}

func NewFlag(o *Flag, defaultType OutputType, defaultStringFormat, defaulErrorFormat string) *Flag {
	*o = Flag{
		outputType:         defaultType,
		outputStringFormat: defaultStringFormat,
		outputErrorFormat:  defaulErrorFormat,
	}
	return o
}

func (o *Flag) Set(v string) error {
	s := strings.Split(v, "=")
	switch s[0] {
	case string(JsonType), string(TableType), string(YamlType):
		o.outputType = OutputType(s[0])
		return nil
	case string(TextType):
		o.outputType = OutputType(s[0])

		if len(s) == 2 {
			o.outputStringFormat = s[1]
		} else if len(s) >= 3 {
			o.outputStringFormat = s[1]
			o.outputErrorFormat = s[2]
		}

		return nil
	default:
		return fmt.Errorf("unsuported output type: %s", s)
	}
}

func (o *Flag) String() string {
	s := string(o.outputType)
	if o.outputType == TextType {
		s = fmt.Sprintf("%s - %s - %s", s, o.outputStringFormat, o.outputErrorFormat)
	}

	return s
}

func (o *Flag) Type() string {
	return string(o.outputType)
}

func (o *Flag) StringFormat() string {
	return o.outputStringFormat
}

func (o *Flag) ErrorFormat() string {
	return o.outputErrorFormat
}

func (o *Flag) Print(w io.Writer, b Formater) error {
	switch o.outputType {
	case JsonType:
		return printJson(w, b.JSON())
	case YamlType:
		return printYaml(w, b.YAML())
	case TableType:
		h, d := b.Table()
		return printTable(w, h, d)
	default: // TextType
		return printText(w, b.Text(o.outputStringFormat, o.outputErrorFormat))
	}
}
