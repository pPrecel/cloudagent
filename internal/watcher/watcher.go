package watcher

import (
	"context"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/internal/gardener"
	"github.com/pPrecel/cloudagent/internal/system"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Context    context.Context
	Logger     *logrus.Entry
	Cache      *agent.ServerCache
	ConfigPath string
}

type watcher struct {
	getConfig    func(string) (*config.Config, error)
	notifyChange func(string) (*system.Notifier, error)
}

func New() *watcher {
	return &watcher{
		getConfig:    config.Read,
		notifyChange: system.NotifyChange,
	}
}

func (w *watcher) Start(o *Options) error {
	o.Logger.Debug("starting cached watcher")
	watcher, err := w.newWatcher(o)
	if err != nil {
		return err
	}
	defer watcher.Stop()
	watcher.Start()

	o.Logger.Info("starting config notifier")
	n, err := w.notifyChange(o.ConfigPath)
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

func (w *watcher) newWatcher(o *Options) (*agent.Watcher, error) {
	o.Logger.Infof("reading config from path: '%s'", o.ConfigPath)
	config, err := w.getConfig(o.ConfigPath)
	if err != nil {
		return nil, err
	}

	o.Logger.Infof("starting state watcher with spec: '%s'", config.PersistentSpec)
	return agent.NewWatcher(agent.WatcherOptions{
		Spec:    config.PersistentSpec,
		Context: o.Context,
		Logger:  o.Logger,
	}, parseWatcherFns(o.Logger, o.Cache.GardenerCache, config)...)
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
