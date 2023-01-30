package watcher

import (
	"context"

	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
)

type cached struct {
	logger     *logrus.Entry
	configPath string
	cache      agent.ResourceGetter

	getConfig func(string) (*config.Config, error)
}

func newCached(cache agent.ResourceGetter, logger *logrus.Entry, configPath string) *cached {
	return &cached{
		logger:     logger,
		configPath: configPath,
		cache:      cache,
		getConfig:  config.Read,
	}
}

func (w *cached) start(ctx context.Context) error {
	w.logger.Debug("starting cached watcher")
	watcher, err := w.newWatcher(ctx)
	if err != nil {
		return err
	}
	defer watcher.Stop()
	watcher.Start()

	<-ctx.Done()
	return nil
}

func (w *cached) newWatcher(ctx context.Context) (*agent.Watcher, error) {
	w.logger.Infof("reading config from path: '%s'", w.configPath)
	config, err := w.getConfig(w.configPath)
	if err != nil {
		return nil, err
	}

	w.logger.Infof("starting state watcher with spec: '%s'", config.PersistentSpec)
	return agent.NewWatcher(agent.WatcherOptions{
		Spec:    config.PersistentSpec,
		Context: ctx,
		Logger:  w.logger,
	}, parseWatcherFns(
		w.logger,
		w.cache.GetGardenerCache(),
		config)...)
}
