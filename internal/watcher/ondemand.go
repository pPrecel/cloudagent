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
	context    context.Context
	logger     *logrus.Entry
	configPath string

	getConfig func(string) (*config.Config, error)
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
	cache := agent.NewCache[*v1beta1.ShootList]()

	cfg, err := rw.getConfig(rw.configPath)
	if err != nil {
		return cache
	}

	fns := parseWatcherFns(rw.logger, cache, cfg)

	for i := range fns {
		fns[i](rw.context)
	}

	return cache
}

func (rw *onDemandWatcher) GetGeneralError() error {
	// reactive watcher never returns general error
	return nil
}
