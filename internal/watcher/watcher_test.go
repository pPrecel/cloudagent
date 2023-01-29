package watcher

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	fixOndemandConfig = func(s string) (*config.Config, error) {
		return &config.Config{
			PersistentSpec: "on-demand",
		}, nil
	}

	fixConfig = func(s string) (*config.Config, error) {
		return &config.Config{
			PersistentSpec: "1s",
		}, nil
	}
)

func TestNewForConfig(t *testing.T) {
	t.Run("read config error", func(t *testing.T) {
		rg, err := NewForConfig(&Options{
			ConfigPath: "",
		})
		assert.Error(t, err)
		assert.Nil(t, rg)
	})
}

func Test_newForConfig(t *testing.T) {
	l := logrus.New()
	l.Out = io.Discard

	t.Run("new on demand", func(t *testing.T) {
		rg, err := newForConfig(&Options{}, fixOndemandConfig)
		assert.NoError(t, err)
		assert.NotNil(t, rg)
	})

	t.Run("watcher context done", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		rg, err := newForConfig(&Options{
			Context: ctx,
			Logger:  l.WithField("test", "test"),
		}, fixConfig)
		assert.NoError(t, err)
		assert.NotNil(t, rg)
	})
}

func Test_setupWatcher(t *testing.T) {
	l := logrus.New()
	l.Out = io.Discard

	t.Run("lib watcher error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		c := &agent.ServerCache{
			GardenerCache: agent.NewCache[*v1beta1.ShootList](),
		}

		setupWatcher(c, &Options{
			Context: ctx,
			Logger:  l.WithField("test", "test"),
		})

		assert.Len(t, c.GardenerCache.Resources(), 0)
	})
}
