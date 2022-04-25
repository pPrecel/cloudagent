package agent

import (
	"context"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	counterFn = func(c chan int, m *sync.Mutex) WatchFn {
		return func(ctx context.Context) {
			m.Lock()
			defer m.Unlock()
			c <- 1
		}
	}
)

func TestNewWatcher(t *testing.T) {
	t.Run("empty watcher", func(t *testing.T) {
		c, e := NewWatcher(WatcherOptions{
			Context: context.Background(),
			Logger:  logrus.New(),
		})
		assert.NoError(t, e)
		assert.Equal(t, 0, len(c.Entries()))
	})

	t.Run("run two funcs", func(t *testing.T) {
		counter := make(chan int)
		m := sync.Mutex{}
		c, e := NewWatcher(WatcherOptions{
			Context: context.Background(),
			Logger:  logrus.New(),
			Spec:    "@every 2s",
		}, counterFn(counter, &m), counterFn(counter, &m))
		defer c.Stop()

		assert.NoError(t, e)

		assert.Equal(t, 2, len(c.Entries()))

		c.Start()

		<-counter
		<-counter
	})

	t.Run("nil fn", func(t *testing.T) {
		c, e := NewWatcher(WatcherOptions{
			Context: context.Background(),
			Logger:  logrus.New(),
			Spec:    "@every 15m",
		}, nil)
		assert.NoError(t, e)
		assert.Equal(t, 1, len(c.Entries()))
	})

	t.Run("wrong spec", func(t *testing.T) {
		c, e := NewWatcher(WatcherOptions{
			Context: context.Background(),
			Logger:  logrus.New(),
		}, nil)
		assert.Error(t, e)
		assert.Nil(t, c)
	})
}
