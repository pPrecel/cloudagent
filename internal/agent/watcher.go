package agent

import (
	"context"

	"github.com/pPrecel/gardener-agent/internal/gardener"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultSpec = "@every 15m"

type WatcherOption struct {
	KubeconfigPath string
	Namespace      string
	Spec           string
	Context        context.Context
	StateSetter    StateSetter
	Logger         *logrus.Logger
}

func NewWatcher(opts WatcherOption) (*cron.Cron, error) {
	opts.Logger.Debug("creating cluster config")
	cfg, err := gardener.NewClusterConfig(opts.KubeconfigPath)
	if err != nil {
		return nil, err
	}

	opts.Logger.Debug("creating gardener client")
	c, err := gardener.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	client := c.Shoots(opts.Namespace)

	cron := cron.New()
	spec := defaultSpec
	if opts.Spec != "" {
		spec = opts.Spec
	}

	context := context.Background()
	if opts.Context != nil {
		context = opts.Context
	}

	opts.Logger.Debug("setting up watchers func")
	_, err = cron.AddFunc(spec, func() {
		opts.Logger.Debug("watching for resources")
		l, err := client.List(context, v1.ListOptions{})
		opts.StateSetter.Set(l)
		if err != nil {
			opts.Logger.Errorf("when watching for shoots: %s", err.Error())
			return
		}

		opts.Logger.Debugf("found %v shoots", len(l.Items))
	})

	return cron, err
}
