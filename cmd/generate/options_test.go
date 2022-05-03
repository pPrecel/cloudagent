package generate

import (
	"runtime"
	"testing"

	command "github.com/pPrecel/cloud-agent/cmd"
	"github.com/stretchr/testify/assert"
)

func Test_options_validate(t *testing.T) {
	type fields struct {
		Options    *command.Options
		ConfigPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "validate",
			fields: fields{
				ConfigPath: func() string {
					_, filename, _, ok := runtime.Caller(0)
					assert.True(t, ok)
					return filename
				}(),
			},
			wantErr: false,
		},
		{
			name: "empty ConfigPath",
			fields: fields{
				ConfigPath: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &options{
				Options:    tt.fields.Options,
				configPath: tt.fields.ConfigPath,
			}
			if err := o.validate(); (err != nil) != tt.wantErr {
				t.Errorf("options.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
