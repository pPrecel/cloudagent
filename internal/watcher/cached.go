package watcher

import (
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/internal/gardener"
	"github.com/pPrecel/cloudagent/internal/system"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
)

type watcher struct {
	options *Options
	cache   *agent.ServerCache

	getConfig    func(string) (*config.Config, error)
	notifyChange func(string) (*system.Notifier, error)
}

func newCached(cache *agent.ServerCache, o *Options) *watcher {
	return &watcher{
		options:      o,
		cache:        cache,
		getConfig:    config.Read,
		notifyChange: system.NotifyChange,
	}
}

func (w *watcher) start() error {
	w.options.Logger.Debug("starting cached watcher")
	watcher, err := w.newWatcher()
	if err != nil {
		return err
	}
	defer watcher.Stop()
	watcher.Start()

	w.options.Logger.Info("starting config notifier")
	n, err := w.notifyChange(w.options.ConfigPath)
	if err != nil {
		return err
	}
	defer n.Stop()

	select {
	case err := <-n.Errors:
		return err
	case <-n.IsMotified:
		return nil
	}
}

func (w *watcher) newWatcher() (*agent.Watcher, error) {
	w.options.Logger.Infof("reading config from path: '%s'", w.options.ConfigPath)
	config, err := w.getConfig(w.options.ConfigPath)
	if err != nil {
		return nil, err
	}

	w.options.Logger.Infof("starting state watcher with spec: '%s'", config.PersistentSpec)
	return agent.NewWatcher(agent.WatcherOptions{
		Spec:    config.PersistentSpec,
		Context: w.options.Context,
		Logger:  w.options.Logger,
	}, parseWatcherFns(w.options.Logger, w.cache.GardenerCache, config)...)
}

func parseWatcherFns(l *logrus.Entry, gardenerCache agent.Cache[*v1beta1.ShootList], config *config.Config) []agent.WatchFn {
	funcs := []agent.WatchFn{}
	for i := range config.GardenerProjects {
		p := config.GardenerProjects[i]
		r := gardenerCache.Register(p.Namespace)

		l.Debugf("creating watcher func for namespace: '%s'", p.Namespace)
		l := l.WithFields(
			logrus.Fields{
				"provider": "gardener",
				"project":  p.Namespace,
			},
		)
		funcs = append(funcs,
			gardener.NewWatchFunc(l, r, p.Namespace, p.KubeconfigPath),
		)
	}

	return funcs
}
