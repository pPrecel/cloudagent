package watcher

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/internal/system"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Context    context.Context
	Logger     *logrus.Entry
	ConfigPath string
}

type watcher struct {
	mu sync.Mutex
	w  agent.ResourceGetter

	o *Options

	getConfig    func(string) (*config.Config, error)
	notifyChange func(string) (*system.Notifier, error)
}

func New(o *Options) *watcher {
	return &watcher{
		o:            o,
		getConfig:    config.Read,
		notifyChange: system.NotifyChange,
	}
}

func (w *watcher) GetGardenerCache() agent.Cache[*v1beta1.ShootList] {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.w != nil {
		return w.w.GetGardenerCache()
	}
	return nil
}

func (w *watcher) GetGeneralError() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.w != nil {
		return w.w.GetGeneralError()
	}
	return errors.New("watcher not implemented")
}

func (w *watcher) Start() {
	for {
		select {
		case <-w.o.Context.Done():
			w.o.Logger.Warn("watcher context done. Exiting")
			return
		default:
			err := w.start()
			if err != nil {
				w.o.Logger.Warnf("runned watcher error: %s", err)
			}

			// wait 1sec to avoid CPU throttling
			time.Sleep(time.Second * 1)
		}
	}
}

type watcherStarter interface {
	start(context.Context) error
}

func (w *watcher) start() error {
	var err error
	starter, err := w.buildWatcher()
	if err != nil {
		return fmt.Errorf("failed to build watcher: %f", err)
	}

	ctx, cancel := context.WithCancel(w.o.Context)
	defer cancel()

	go starter.start(ctx)

	w.o.Logger.Info("starting config notifier")
	n, err := w.notifyChange(w.o.ConfigPath)
	if err != nil {
		return err
	}
	defer n.Stop()

	select {
	case err := <-n.Errors:
		return err
	case <-n.IsMotified:
		w.o.Logger.Info("configuration midyfication detected")
		return nil
	}
}

func (w *watcher) buildWatcher() (watcherStarter, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	cfg, err := w.getConfig(w.o.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %s", err)
	}

	if cfg.PersistentSpec == "on-demand" {
		ondemand := newOnDemand(w.o)
		w.w = ondemand
		return ondemand, nil
	}

	w.w = newCache(w.w)
	return newCached(w.w, w.o.Logger, w.o.ConfigPath), nil
}

func newCache(actualCache agent.ResourceGetter) agent.ResourceGetter {
	if actualCache == nil {
		return &agent.ServerCache{
			GardenerCache: agent.NewCache[*v1beta1.ShootList](),
		}
	}

	switch actualCache.(type) {
	case *agent.ServerCache:
		c := actualCache.(*agent.ServerCache)
		c.GardenerCache.Clean()
		c.GeneralError = nil
		return c
	default:
		return newCache(nil)
	}
}
