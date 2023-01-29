package watcher

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testShootList = &v1beta1.ShootList{
		Items: []v1beta1.Shoot{
			{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-shoot",
				},
			},
		},
	}
)

func TestNewOnDemand(t *testing.T) {
	l := &logrus.Entry{
		Logger: logrus.New(),
	}
	l.Logger.Out = io.Discard

	t.Run("build new watcher", func(t *testing.T) {
		assert.NotNil(t, newOnDemand(&Options{
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
			parseWatcherFns: func(e *logrus.Entry, a agent.Cache[*v1beta1.ShootList], c *config.Config) []agent.WatchFn {
				return []agent.WatchFn{}
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
	t.Run("update cache using fns", func(t *testing.T) {
		cache := agent.NewCache[*v1beta1.ShootList]()
		cache.Register("test-1").Set(nil, nil)

		w := onDemandWatcher{
			cache: cache,
			fns: []agent.WatchFn{
				func(ctx context.Context) {
					cache.Register("test-1").Set(testShootList, nil)
				},
			},
		}

		cache = w.GetGardenerCache()

		assert.Equal(t, testShootList, cache.Resources()["test-1"].Get().Value)
		assert.NoError(t, cache.Resources()["test-1"].Get().Error)
	})
}
