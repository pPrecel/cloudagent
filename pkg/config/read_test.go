package config

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testdataConfig = &Config{
		PersistentSpec: "@every 2s",
		GardenerProjects: []GardenerProject{
			{
				Namespace:      "n1",
				KubeconfigPath: "k1",
			},
			{
				Namespace:      "n2",
				KubeconfigPath: "k2",
			},
			{
				Namespace:      "n3",
				KubeconfigPath: "k3",
			},
		},
		GCPProjects: []GCPProject{},
	}
)

func TestGetConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "read config",
			args: args{
				path: func() string {
					_, filename, _, ok := runtime.Caller(0)
					assert.True(t, ok)

					return filepath.Join(filepath.Dir(filename), "/testdata/test.conf.yaml")
				}(),
			},
			want:    testdataConfig,
			wantErr: false,
		},
		{
			name: "empty config",
			args: args{
				path: func() string {
					_, filename, _, ok := runtime.Caller(0)
					assert.True(t, ok)

					return filepath.Join(filepath.Dir(filename), "/testdata/empty.conf.yaml")
				}(),
			},
			want:    &Config{},
			wantErr: false,
		},
		{
			name: "config does not exist",
			args: args{
				path: "/this/path/does/not/exist",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "file is not a config",
			args: args{
				path: func() string {
					_, filename, _, ok := runtime.Caller(0)
					assert.True(t, ok)

					return filepath.Join(filepath.Dir(filename), "/testdata/file.txt")
				}(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConfig(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
