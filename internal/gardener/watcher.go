package gardener

import (
	"context"

	"github.com/pPrecel/cloud-agent/internal/agent"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WatchOptions struct {
	KubeconfigPath string
	Namespace      string
	StateSetter    StateSetter
	Logger         *logrus.Logger
}

func NewWatchFunc(opts WatchOptions) (agent.WatchFn, error) {
	opts.Logger.Debug("creating cluster config")
	cfg, err := newClusterConfig(opts.KubeconfigPath)
	if err != nil {
		return nil, err
	}

	opts.Logger.Debug("creating gardener client")
	c, err := newClient(cfg)
	if err != nil {
		return nil, err
	}
	client := c.Shoots(opts.Namespace)

	opts.Logger.Debug("setting up watchers func")
	return func(context context.Context) {
		opts.Logger.Debug("watching for resources")
		l, err := client.List(context, v1.ListOptions{})
		opts.StateSetter.Set(l)
		if err != nil {
			opts.Logger.Errorf("when watching for shoots: %s", err.Error())
			return
		}

		opts.Logger.Debugf("found %v shoots", len(l.Items))
	}, nil
}
