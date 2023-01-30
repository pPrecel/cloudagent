package watcher

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/internal/system"
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

func TestNew(t *testing.T) {
	t.Run("read config error", func(t *testing.T) {
		rg := New(&Options{
			ConfigPath: "",
		})
		assert.NotNil(t, rg)
	})
}

type resourceGetterStub struct {
	ch    chan struct{}
	cache agent.Cache[*v1beta1.ShootList]
	err   error
}

func (stub *resourceGetterStub) GetGardenerCache() agent.Cache[*v1beta1.ShootList] {
	<-stub.ch
	return stub.cache
}

func (stub *resourceGetterStub) GetGeneralError() error {
	<-stub.ch
	return stub.err
}

func TestWatcher_GetGardenerCache(t *testing.T) {
	t.Run("lock and return cache", func(t *testing.T) {
		ch := make(chan struct{})
		resourceGetter := resourceGetterStub{
			ch: ch,
		}
		w := watcher{
			w: &resourceGetter,
		}

		stop := make(chan struct{})
		go func() {
			c := w.GetGardenerCache()
			assert.Nil(t, c)
			stop <- struct{}{}
		}()

		ch <- struct{}{}
		<-stop
	})
	t.Run("nil cache", func(t *testing.T) {
		w := watcher{
			w: nil,
		}

		assert.Nil(t, w.GetGardenerCache())
	})
}

func Test_watcher_GetGeneralError(t *testing.T) {
	t.Run("lock and return cache", func(t *testing.T) {
		ch := make(chan struct{})
		resourceGetter := resourceGetterStub{
			ch: ch,
		}
		w := watcher{
			w: &resourceGetter,
		}

		stop := make(chan struct{})
		go func() {
			c := w.GetGeneralError()
			assert.Nil(t, c)
			stop <- struct{}{}
		}()

		ch <- struct{}{}
		<-stop
	})
	t.Run("nil cache", func(t *testing.T) {
		w := watcher{
			w: nil,
		}

		assert.Error(t, w.GetGeneralError())
	})
}

func Test_watcher_Start(t *testing.T) {
	l := &logrus.Entry{
		Logger: logrus.New(),
	}
	l.Logger.Out = io.Discard

	t.Run("context done", func(t *testing.T) {
		w := watcher{
			o: &Options{
				Context: fixCanceledContext(),
				Logger:  l,
			},
		}

		w.Start()
	})

	t.Run("do not panic on error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		stop := make(chan struct{})
		w := watcher{
			getConfig: func(s string) (*config.Config, error) {
				stop <- struct{}{}
				return nil, errors.New("test error")
			},
			o: &Options{
				Context: ctx,
				Logger:  l,
			},
		}

		go w.Start()

		<-stop
		cancel()
	})
}

func Test_watcher_start(t *testing.T) {
	l := &logrus.Entry{
		Logger: logrus.New(),
	}
	l.Logger.Out = io.Discard

	t.Run("start on-demand", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		modified := make(chan interface{})
		notifier := &system.Notifier{
			IsMotified: modified,
			Stop:       func() {},
		}

		stop := make(chan interface{})
		w := watcher{
			getConfig: func(s string) (*config.Config, error) {
				return &config.Config{
					PersistentSpec: "on-demand",
				}, nil
			},
			notifyChange: func(s string) (*system.Notifier, error) {
				return notifier, nil
			},
			o: &Options{
				Context: ctx,
				Logger:  l,
			},
		}

		go func() {
			err := w.start()
			assert.Nil(t, err)
			stop <- struct{}{}
		}()

		modified <- struct{}{}
		<-stop
	})

	t.Run("start cached", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		modified := make(chan interface{})
		notifier := &system.Notifier{
			IsMotified: modified,
			Stop:       func() {},
		}

		stop := make(chan interface{})
		w := watcher{
			getConfig: func(s string) (*config.Config, error) {
				return &config.Config{
					PersistentSpec: "@every 120s",
				}, nil
			},
			notifyChange: func(s string) (*system.Notifier, error) {
				return notifier, nil
			},
			o: &Options{
				Context: ctx,
				Logger:  l,
			},
		}

		go func() {
			err := w.start()
			assert.Nil(t, err)
			stop <- struct{}{}
		}()

		modified <- struct{}{}
		<-stop
	})

	t.Run("start cached when cache is not nil", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		modified := make(chan interface{})
		notifier := &system.Notifier{
			IsMotified: modified,
			Stop:       func() {},
		}

		stop := make(chan interface{})
		w := watcher{
			w: &agent.ServerCache{
				GardenerCache: agent.NewCache[*v1beta1.ShootList](),
			},
			getConfig: func(s string) (*config.Config, error) {
				return &config.Config{
					PersistentSpec: "@every 120s",
				}, nil
			},
			notifyChange: func(s string) (*system.Notifier, error) {
				return notifier, nil
			},
			o: &Options{
				Context: ctx,
				Logger:  l,
			},
		}

		go func() {
			err := w.start()
			assert.Nil(t, err)
			stop <- struct{}{}
		}()

		modified <- struct{}{}
		<-stop
	})

	t.Run("start cached when cache is on-demand watcher", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		modified := make(chan interface{})
		notifier := &system.Notifier{
			IsMotified: modified,
			Stop:       func() {},
		}

		stop := make(chan interface{})
		w := watcher{
			w: &ondemand{},
			getConfig: func(s string) (*config.Config, error) {
				return &config.Config{
					PersistentSpec: "@every 120s",
				}, nil
			},
			notifyChange: func(s string) (*system.Notifier, error) {
				return notifier, nil
			},
			o: &Options{
				Context: ctx,
				Logger:  l,
			},
		}

		go func() {
			err := w.start()
			assert.Nil(t, err)
			stop <- struct{}{}
		}()

		modified <- struct{}{}
		<-stop
	})

	t.Run("handle notify error channel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		chanErr := make(chan error)
		notifier := &system.Notifier{
			Errors: chanErr,
			Stop:   func() {},
		}

		stop := make(chan interface{})
		w := watcher{
			getConfig: func(s string) (*config.Config, error) {
				return &config.Config{
					PersistentSpec: "on-demand",
				}, nil
			},
			notifyChange: func(s string) (*system.Notifier, error) {
				return notifier, nil
			},
			o: &Options{
				Context: ctx,
				Logger:  l,
			},
		}

		go func() {
			err := w.start()
			assert.Error(t, err)
			stop <- struct{}{}
		}()

		chanErr <- errors.New("test error")
		<-stop
	})

	t.Run("notify error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		w := watcher{
			getConfig: func(s string) (*config.Config, error) {
				return &config.Config{
					PersistentSpec: "on-demand",
				}, nil
			},
			notifyChange: func(s string) (*system.Notifier, error) {
				return nil, errors.New("test error")
			},
			o: &Options{
				Context: ctx,
				Logger:  l,
			},
		}

		err := w.start()
		assert.Error(t, err)
	})
}

func fixCanceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
