package gardener

import (
	"context"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

//go:generate mockery --name=Client --output=automock --outpkg=automock
type Client interface {
	List(context.Context, v1.ListOptions) (*v1beta1.ShootList, error)
}

func NewWatchFunc(l *logrus.Entry, r agent.RegisteredResource[*v1beta1.ShootList], namespace, kubeconfig string) agent.WatchFn {
	return newWatchFunc(l, r, newClientBuilder(l, newClusterConfig, namespace, kubeconfig))
}

func newWatchFunc(l *logrus.Entry, r agent.RegisteredResource[*v1beta1.ShootList], clientBuilder func() (Client, error)) agent.WatchFn {
	l.Debug("setting up watchers func")
	var c Client
	var err error

	return func(context context.Context) {
		l.Debug("watching for resources")
		if c == nil || err != nil {
			l.Info("building new gardener client")
			c, err = clientBuilder()
			if err != nil {
				l.Errorf("when creating gardener client: %s", err.Error())
				r.Set(nil, err)
				return
			}
		}

		list, err := c.List(context, v1.ListOptions{})
		r.Set(list, err)
		if err != nil {
			l.Errorf("when watching for shoots: %s", err.Error())
			return
		}

		l.Debugf("found %v shoots", len(list.Items))
	}
}

func newClientBuilder(l *logrus.Entry, buildConfig func(string) (*rest.Config, error), namespace, kubeconfig string) func() (Client, error) {
	return func() (Client, error) {
		l.Debugf("creating cluster config for kubeconfig: %s", kubeconfig)
		cfg, err := buildConfig(kubeconfig)
		if err != nil {
			return nil, err
		}

		l.Debug("creating gardener client")
		c, err := newClient(cfg)
		if err != nil {
			return nil, err
		}

		return c.Shoots(namespace), nil
	}
}
