package watcher

import (
	"context"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
)

type onDemandWatcher struct {
	cache agent.Cache[*v1beta1.ShootList]
	fns   []agent.WatchFn

	context    context.Context
	logger     *logrus.Entry
	configPath string

	getConfig       func(string) (*config.Config, error)
	parseWatcherFns func(*logrus.Entry, agent.Cache[*v1beta1.ShootList], *config.Config) []agent.WatchFn
}

func newOnDemand(o *Options) *onDemandWatcher {
	return &onDemandWatcher{
		context:         o.Context,
		logger:          o.Logger,
		cache:           agent.NewCache[*v1beta1.ShootList](),
		configPath:      o.ConfigPath,
		getConfig:       config.Read,
		parseWatcherFns: parseWatcherFns,
	}
}

func (rw *onDemandWatcher) GetGardenerCache() agent.Cache[*v1beta1.ShootList] {
	rw.cache.Clean()
	for i := range rw.fns {
		rw.fns[i](rw.context)
	}

	return rw.cache
}

func (rw *onDemandWatcher) GetGeneralError() error {
	cfg, err := rw.getConfig(rw.configPath)
	if err != nil {
		return err
	}

	rw.fns = rw.parseWatcherFns(rw.logger, rw.cache, cfg)

	return nil
}
