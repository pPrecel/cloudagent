package watcher

import (
	"context"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
)

type ondemand struct {
	cache agent.Cache[*v1beta1.ShootList]
	fns   []agent.WatchFn

	context    context.Context
	logger     *logrus.Entry
	configPath string

	getConfig       func(string) (*config.Config, error)
	parseWatcherFns func(*logrus.Entry, agent.Cache[*v1beta1.ShootList], *config.Config) []agent.WatchFn
}

func newOnDemand(o *Options) *ondemand {
	return &ondemand{
		context:         o.Context,
		logger:          o.Logger,
		cache:           agent.NewCache[*v1beta1.ShootList](),
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

func (w *ondemand) GetGardenerCache() agent.Cache[*v1beta1.ShootList] {
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
