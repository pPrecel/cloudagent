package watcher

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/pPrecel/cloudagent/pkg/cache"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	fixConfigFn = func(s string) (*config.Config, error) {
		return &config.Config{
			PersistentSpec: "@every 2m",
			GardenerProjects: []config.GardenerProject{
				{
					Namespace:      "test-namespace",
					KubeconfigPath: "path",
				},
			},
		}, nil
	}
)

func TestNewCached(t *testing.T) {
	t.Run("new cached watcher", func(t *testing.T) {
		assert.NotNil(t, newCached(nil, nil, ""))
	})
}

func Test_cached_Start(t *testing.T) {
	l := &logrus.Entry{
		Logger: logrus.New(),
	}
	l.Logger.Out = io.Discard

	type fields struct {
		logger     *logrus.Entry
		configPath string
		getConfig  func(string) (*config.Config, error)
	}
	type args struct {
		context context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		cache   *cache.ServerCache
		wantErr bool
	}{
		{
			name: "context done",
			fields: fields{
				configPath: "",
				logger:     l,
				getConfig: func(s string) (*config.Config, error) {
					return &config.Config{
						PersistentSpec: "@every 120s",
						GardenerProjects: []config.GardenerProject{
							{
								Namespace:      "test",
								KubeconfigPath: "/test/path",
							},
						},
					}, nil
				},
			},
			args: args{
				context: fixCanceledContext(),
			},
			cache: &cache.ServerCache{
				GardenerCache: cache.NewGardenerCache(),
			},
			wantErr: false,
		},
		{
			name: "getConfig error",
			fields: fields{
				configPath: "",
				logger:     l,
				getConfig: func(s string) (*config.Config, error) {
					return nil, errors.New("test error")
				},
			},
			args: args{
				context: context.Background(),
			},
			cache: &cache.ServerCache{
				GardenerCache: cache.NewGardenerCache(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &cached{
				logger:     tt.fields.logger,
				configPath: tt.fields.configPath,
				getConfig:  tt.fields.getConfig,
				cache:      tt.cache,
			}

			if err := w.start(tt.args.context); (err != nil) != tt.wantErr {
				t.Errorf("watcher.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
