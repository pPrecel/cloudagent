package formater

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pPrecel/cloudagent/pkg/config"
)

var (
	testConfig = &config.Config{
		PersistentSpec: "@every 2s",
		GardenerProjects: []config.GardenerProject{
			{
				Namespace:      "namespace1",
				KubeconfigPath: "/any/path.yaml",
			},
			{},
		},
		GCPProjects: []config.GCPProject{
			// TODO: add me pls
		},
	}

	testConfigRows = [][]string{
		{"namespace1", "/any/path.yaml"},
		{"", ""},
	}
)

func Test_configFormater_YAML(t *testing.T) {
	type fields struct {
		err error
		cfg *config.Config
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "nil arguments",
			fields: fields{
				err: nil,
				cfg: nil,
			},
			want: map[string]interface{}{},
		},
		{
			name: "empty arguments",
			fields: fields{
				err: nil,
				cfg: &config.Config{},
			},
			want: &config.Config{},
		},
		{
			name: "error",
			fields: fields{
				err: errors.New("test error"),
				cfg: nil,
			},
			want: map[string]interface{}{},
		},
		{
			name: "get json",
			fields: fields{
				err: nil,
				cfg: testConfig,
			},
			want: testConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewConfig(tt.fields.err, tt.fields.cfg)
			if got := f.YAML(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configFormater.YAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configFormater_JSON(t *testing.T) {
	type fields struct {
		err error
		cfg *config.Config
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "nil arguments",
			fields: fields{
				err: nil,
				cfg: nil,
			},
			want: map[string]interface{}{},
		},
		{
			name: "empty arguments",
			fields: fields{
				err: nil,
				cfg: &config.Config{},
			},
			want: &config.Config{},
		},
		{
			name: "error",
			fields: fields{
				err: errors.New("test error"),
				cfg: nil,
			},
			want: map[string]interface{}{},
		},
		{
			name: "get yaml",
			fields: fields{
				err: nil,
				cfg: testConfig,
			},
			want: testConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewConfig(tt.fields.err, tt.fields.cfg)
			if got := f.JSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configFormater.JSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configFormater_Table(t *testing.T) {
	type fields struct {
		err error
		cfg *config.Config
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
		want1  [][]string
	}{
		{
			name: "nil arguments",
			fields: fields{
				err: nil,
				cfg: nil,
			},
			want:  configHeaders,
			want1: [][]string{},
		},
		{
			name: "empty arguments",
			fields: fields{
				err: nil,
				cfg: &config.Config{},
			},
			want:  configHeaders,
			want1: [][]string{},
		},
		{
			name: "error",
			fields: fields{
				err: errors.New("test error"),
				cfg: nil,
			},
			want:  configHeaders,
			want1: [][]string{},
		},
		{
			name: "get table",
			fields: fields{
				err: nil,
				cfg: testConfig,
			},
			want:  configHeaders,
			want1: testConfigRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewConfig(tt.fields.err, tt.fields.cfg)
			got, got1 := f.Table()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configFormater.Table() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("configFormater.Table() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_configFormater_Text(t *testing.T) {
	type fields struct {
		err error
		cfg *config.Config
	}
	type args struct {
		outFormat string
		errFormat string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "nil arguments",
			fields: fields{
				err: nil,
				cfg: nil,
			},
			args: args{
				outFormat: "",
				errFormat: "$e",
			},
			want: "nil config",
		},
		{
			name: "empty config",
			fields: fields{
				err: nil,
				cfg: &config.Config{},
			},
			args: args{
				outFormat: "$a$g$G$p",
				errFormat: "$e",
			},
			want: "000",
		},
		{
			name: "error",
			fields: fields{
				err: errors.New("test error"),
				cfg: nil,
			},
			args: args{
				outFormat: "$a$g$G$p",
				errFormat: "$e",
			},
			want: "test error",
		},
		{
			name: "get text",
			fields: fields{
				err: nil,
				cfg: testConfig,
			},
			args: args{
				outFormat: "$a/$g/$G/$p",
				errFormat: "$e",
			},
			want: "2/2/0/@every 2s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewConfig(tt.fields.err, tt.fields.cfg)
			if got := f.Text(tt.args.outFormat, tt.args.errFormat); got != tt.want {
				t.Errorf("configFormater.Text() = %v, want %v", got, tt.want)
			}
		})
	}
}
