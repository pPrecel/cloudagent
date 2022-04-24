package serve

import (
	"runtime"
	"testing"

	"github.com/pPrecel/cloud-agent/internal/command"
	"github.com/stretchr/testify/assert"
)

func Test_options_validate(t *testing.T) {
	type fields struct {
		Options        *command.Options
		KubeconfigPath string
		Namespace      string
		CronSpec       string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "validate",
			fields: fields{
				KubeconfigPath: func() string {
					_, filename, _, ok := runtime.Caller(0)
					assert.True(t, ok)
					return filename
				}(),
				Namespace: "anything",
				CronSpec:  "2s",
			},
			wantErr: false,
		},
		{
			name: "empty namespace",
			fields: fields{
				KubeconfigPath: func() string {
					_, filename, _, ok := runtime.Caller(0)
					assert.True(t, ok)
					return filename
				}(),
			},
			wantErr: true,
		},
		{
			name: "empty kubeconfigPath",
			fields: fields{
				KubeconfigPath: "",
				Namespace:      "anything",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &options{
				Options:        tt.fields.Options,
				kubeconfigPath: tt.fields.KubeconfigPath,
				namespace:      tt.fields.Namespace,
				cronSpec:       tt.fields.CronSpec,
			}
			if err := o.validate(); (err != nil) != tt.wantErr {
				t.Errorf("options.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
