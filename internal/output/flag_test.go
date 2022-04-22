package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutput_Set(t *testing.T) {
	type outputFields struct {
		outputType   OutputType
		stringFormat string
		errorFormat  string
	}
	tests := []struct {
		name    string
		arg     string
		out     outputFields
		wantErr bool
	}{
		{
			name:    "error with empty arg",
			wantErr: true,
		},
		{
			name: "set value for json type",
			arg:  "json",
			out: outputFields{
				outputType: JsonType,
			},
			wantErr: false,
		},
		{
			name: "set value with too many params for json type",
			arg:  "json=aa=bb=cc",
			out: outputFields{
				outputType: JsonType,
			},
			wantErr: false,
		},
		{
			name: "set value for table type",
			arg:  "table",
			out: outputFields{
				outputType: TableType,
			},
			wantErr: false,
		},
		{
			name: "set value with too many params for table type",
			arg:  "table=aa=bb=cc",
			out: outputFields{
				outputType: TableType,
			},
			wantErr: false,
		},
		{
			name: "set value for text type",
			arg:  "text",
			out: outputFields{
				outputType: TextType,
			},
			wantErr: false,
		},
		{
			name: "set value with one params for table type",
			arg:  "text=aa",
			out: outputFields{
				outputType:   TextType,
				stringFormat: "aa",
			},
			wantErr: false,
		},
		{
			name: "set value with two params for table type",
			arg:  "text=aa=bb",
			out: outputFields{
				outputType:   TextType,
				stringFormat: "aa",
				errorFormat:  "bb",
			},
			wantErr: false,
		},
		{
			name: "set value with a few params for table type",
			arg:  "text=aa=bb=cc",
			out: outputFields{
				outputType:   TextType,
				stringFormat: "aa",
				errorFormat:  "bb",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Output{}
			if err := o.Set(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("Output.Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if o.outputErrorFormat != tt.out.errorFormat ||
				o.outputStringFormat != tt.out.stringFormat ||
				o.outputType != tt.out.outputType {
				t.Errorf("Fields are not the same: %v, want %v", o, tt.out)
			}
		})
	}
}

func TestOutput_SmallerMethods(t *testing.T) {
	t.Run("create table output and change it to text", func(t *testing.T) {
		output := &Output{}
		o := New(output, TableType, "", "")
		assert.Equal(t, output, o)

		assert.NoError(t, o.Set("table"))
		assert.Equal(t, string(TableType), o.Type())
		assert.Equal(t, string(TableType), o.String())
		assert.Equal(t, "", o.ErrorFormat())
		assert.Equal(t, "", o.StringFormat())

		assert.NoError(t, o.Set("text=aa=bb"))
		assert.Equal(t, string(TextType), o.Type())
		assert.Equal(t, "text - aa - bb", o.String())
		assert.Equal(t, "bb", o.ErrorFormat())
		assert.Equal(t, "aa", o.StringFormat())
	})
}
