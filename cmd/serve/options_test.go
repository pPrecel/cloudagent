package serve

import (
	"testing"

	command "github.com/pPrecel/cloudagent/cmd"
)

func Test_options_validate(t *testing.T) {
	type fields struct {
		Options    *command.Options
		configPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "validate",
			fields: fields{
				configPath: "any",
			},
			wantErr: false,
		},
		{
			name: "empty configPath",
			fields: fields{
				configPath: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &options{
				Options:    tt.fields.Options,
				configPath: tt.fields.configPath,
			}
			if err := o.validate(); (err != nil) != tt.wantErr {
				t.Errorf("options.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
