package watcher

import (
	"context"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
)

type NewOnDemandOptions struct {
	Context    context.Context
	Logger     *logrus.Entry
	ConfigPath string
}

type onDemandWatcher struct {
	cache  agent.Cache[*v1beta1.ShootList]
	config *config.Config

	context    context.Context
	logger     *logrus.Entry
	configPath string

	getConfig       func(string) (*config.Config, error)
	parseWatcherFns func(*logrus.Entry, agent.Cache[*v1beta1.ShootList], *config.Config) []agent.WatchFn
}

func NewOnDemand(o *NewOnDemandOptions) *onDemandWatcher {
	o.Logger.Info("starting on demand watcher")
	return &onDemandWatcher{
		context:    o.Context,
		logger:     o.Logger,
		configPath: o.ConfigPath,
		getConfig:  config.Read,
	}
}

func (rw *onDemandWatcher) GetGardenerCache() agent.Cache[*v1beta1.ShootList] {
	fns := rw.parseWatcherFns(rw.logger, rw.cache, rw.config)

	for i := range fns {
		fns[i](rw.context)
	}

	return rw.cache
}

func (rw *onDemandWatcher) GetGeneralError() error {
	rw.cache = agent.NewCache[*v1beta1.ShootList]()

	var err error
	rw.config, err = rw.getConfig(rw.configPath)
	return err
}
