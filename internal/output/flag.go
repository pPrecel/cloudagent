package output

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

type OutputType string

const (
	JsonType  OutputType = "json"
	TableType OutputType = "table"
	TextType  OutputType = "text"
)

var _ pflag.Value = &Output{}

type Output struct {
	outputType         OutputType
	outputStringFormat string
	outputErrorFormat  string
}

func New(o *Output, defaultType OutputType, defaultStringFormat, defaulErrorFormat string) *Output {
	*o = Output{
		outputType:         defaultType,
		outputStringFormat: defaultStringFormat,
		outputErrorFormat:  defaulErrorFormat,
	}
	return o
}

func (o *Output) Set(v string) error {
	s := strings.Split(v, "=")
	switch s[0] {
	case string(JsonType), string(TableType):
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

func (o *Output) String() string {
	s := string(o.outputType)
	if o.outputType == TextType {
		s = fmt.Sprintf("%s - %s - %s", s, o.outputStringFormat, o.outputErrorFormat)
	}

	return s
}

func (o *Output) Type() string {
	return string(o.outputType)
}

func (o *Output) StringFormat() string {
	return o.outputStringFormat
}

func (o *Output) ErrorFormat() string {
	return o.outputErrorFormat
}
