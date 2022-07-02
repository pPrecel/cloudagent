package watcher

import (
	"context"

	v1beta1_apis "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/internal/gardener"
	"github.com/pPrecel/cloudagent/internal/system"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

type Options struct {
	Context    context.Context
	Logger     *logrus.Logger
	Cache      agent.Cache[*v1beta1_apis.ShootList]
	ConfigPath string
}

type watcher struct {
	getConfig        func(string) (*config.Config, error)
	newClusterConfig func(string) (*rest.Config, error)
	notifyChange     func(string) (*system.Notifier, error)
}

func NewWatcher() *watcher {
	return &watcher{
		getConfig:        config.Read,
		newClusterConfig: gardener.NewClusterConfig,
		notifyChange:     system.NotifyChange,
	}
}

func (w *watcher) Start(o *Options) error {
	o.Logger.Debug("starting watcher")
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

	funcs := []agent.WatchFn{}
	for i := range config.GardenerProjects {
		p := config.GardenerProjects[i]

		o.Logger.Debugf("creating cluster config for kubeconfig: %s", p.KubeconfigPath)
		cfg, err := w.newClusterConfig(p.KubeconfigPath)
		if err != nil {
			return nil, err
		}

		o.Logger.Debug("creating gardener client")
		c, err := gardener.NewClient(cfg)
		if err != nil {
			return nil, err
		}

		r := o.Cache.Register(p.Namespace)

		o.Logger.Debugf("creeating watcher func for namespace: '%s'", p.Namespace)
		funcs = append(funcs,
			gardener.NewWatchFunc(o.Logger, c.Shoots(p.Namespace), r),
		)
	}

	return agent.NewWatcher(agent.WatcherOptions{
		Spec:    config.PersistentSpec,
		Context: o.Context,
		Logger:  o.Logger,
	}, funcs...)
}
