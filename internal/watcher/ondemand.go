package watcher

import (
	"context"

	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/cache"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
)

type ondemand struct {
	cache cache.GardenerCache
	fns   []agent.WatchFn

	context    context.Context
	logger     *logrus.Entry
	configPath string

	getConfig       func(string) (*config.Config, error)
	parseWatcherFns func(*logrus.Entry, cache.GardenerCache, *config.Config) []agent.WatchFn
}

func newOnDemand(o *Options) *ondemand {
	return &ondemand{
		context:         o.Context,
		logger:          o.Logger,
		cache:           cache.NewGardenerCache(),
		configPath:      o.ConfigPath,
		getConfig:       config.Read,
		parseWatcherFns: parseWatcherFns,
	}
}

func (w *ondemand) start(ctx context.Context) error {
	// wait for cotext Done only
	<-ctx.Done()
	return nil
}

func (w *ondemand) GetGardenerCache() cache.GardenerCache {
	for i := range w.fns {
		w.fns[i](w.context)
	}

	return w.cache
}

func (w *ondemand) GetGeneralError() error {
	w.cache.Clean()
	cfg, err := w.getConfig(w.configPath)
	if err != nil {
		return err
	}

	w.fns = w.parseWatcherFns(w.logger, w.cache, cfg)

	return nil
}
