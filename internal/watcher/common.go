package watcher

import (
	"github.com/pPrecel/cloudagent/internal/gardener"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/cache"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
)

func parseWatcherFns(l *logrus.Entry, gardenerCache cache.GardenerCache, config *config.Config) []agent.WatchFn {
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
