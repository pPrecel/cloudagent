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
		w, e := NewWatcher(WatcherOptions{
			Context: context.Background(),
			Logger:  logrus.New(),
		})
		assert.NoError(t, e)
		assert.Equal(t, 0, len(w.c.Entries()))
	})

	t.Run("run two funcs", func(t *testing.T) {
		counter := make(chan int)
		m := sync.Mutex{}
		w, e := NewWatcher(WatcherOptions{
			Context: context.Background(),
			Logger:  logrus.New(),
			Spec:    "@every 2s",
		}, counterFn(counter, &m), counterFn(counter, &m))
		defer w.Stop()

		assert.NoError(t, e)

		assert.Equal(t, 2, len(w.c.Entries()))

		w.Start()

		<-counter
		<-counter
	})

	t.Run("nil fn", func(t *testing.T) {
		w, e := NewWatcher(WatcherOptions{
			Context: context.Background(),
			Logger:  logrus.New(),
			Spec:    "@every 15m",
		}, nil)
		assert.NoError(t, e)
		assert.Equal(t, 1, len(w.c.Entries()))
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
