package watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Context    context.Context
	Logger     *logrus.Entry
	ConfigPath string
}

func NewForConfig(o *Options) (agent.ResourceGetter, error) {
	return newForConfig(o, config.Read)
}

func newForConfig(o *Options, getConfig func(string) (*config.Config, error)) (agent.ResourceGetter, error) {
	cfg, err := getConfig(o.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %s", err)
	}

	if cfg.PersistentSpec == "onDemand" {
		return newOnDemand(o), nil
	} else {
		cache := &agent.ServerCache{
			GardenerCache: agent.NewCache[*v1beta1.ShootList](),
		}

		go setupWatcher(cache, o)

		return cache, nil
	}
}

func setupWatcher(cache *agent.ServerCache, o *Options) {
	for {

		select {
		case <-o.Context.Done():
			o.Logger.Warn("watcher context done. Exiting")
			return
		default:
			cache.GeneralError = nil
			startWatcher(cache, o)

			// wait 1sec to avoid CPU throttling
			time.Sleep(time.Second * 1)
		}
	}
}

func startWatcher(cache *agent.ServerCache, o *Options) {
	if err := newCached(cache, &Options{
		Context:    o.Context,
		Logger:     o.Logger,
		ConfigPath: o.ConfigPath,
	}).start(); err != nil {
		o.Logger.Warn(err)
		cache.GeneralError = err
	}

	o.Logger.Info("configuration midyfication detected")

	o.Logger.Info("cleaning up cache")
	cache.GardenerCache.Clean()
}
