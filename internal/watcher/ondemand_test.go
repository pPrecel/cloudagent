package watcher

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewOnDemand(t *testing.T) {
	l := &logrus.Entry{
		Logger: logrus.New(),
	}
	l.Logger.Out = io.Discard

	t.Run("build new watcher", func(t *testing.T) {
		assert.NotNil(t, NewOnDemand(&NewOnDemandOptions{
			Logger: l,
		}))
	})
}

func Test_onDemandWatcher_GetGeneralError(t *testing.T) {
	t.Run("get nil", func(t *testing.T) {
		w := onDemandWatcher{
			getConfig: func(s string) (*config.Config, error) {
				return nil, nil
			},
		}

		assert.Nil(t, w.GetGeneralError())
	})

	t.Run("get general error", func(t *testing.T) {
		w := onDemandWatcher{
			getConfig: func(s string) (*config.Config, error) {
				return nil, errors.New("test error")
			},
		}

		assert.NotNil(t, w.GetGeneralError())
	})
}

func Test_onDemandWatcher_GetGardenerCache(t *testing.T) {
	type fields struct {
		cache           agent.Cache[*v1beta1.ShootList]
		config          *config.Config
		parseWatcherFns func(*logrus.Entry, agent.Cache[*v1beta1.ShootList], *config.Config) []agent.WatchFn
	}
	tests := []struct {
		name   string
		fields fields
		want   agent.Cache[*v1beta1.ShootList]
	}{
		{
			name: "get cache",
			fields: fields{
				cache: agent.NewCache[*v1beta1.ShootList](),
				parseWatcherFns: func(e *logrus.Entry, a agent.Cache[*v1beta1.ShootList], c *config.Config) []agent.WatchFn {
					return []agent.WatchFn{
						func(ctx context.Context) {},
					}
				},
			},
			want: agent.NewCache[*v1beta1.ShootList](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &onDemandWatcher{
				cache:           tt.fields.cache,
				config:          tt.fields.config,
				parseWatcherFns: tt.fields.parseWatcherFns,
			}
			if got := rw.GetGardenerCache(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("onDemandWatcher.GetGardenerCache() = %v, want %v", got, tt.want)
			}
		})
	}
}
