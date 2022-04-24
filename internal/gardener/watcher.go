package gardener

import (
	"context"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloud-agent/internal/agent"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//go:generate mockery --name=Client --output=automock --outpkg=automock
type Client interface {
	List(context.Context, v1.ListOptions) (*v1beta1.ShootList, error)
}

func NewWatchFunc(l *logrus.Logger, c Client, s StateSetter) agent.WatchFn {
	l.Debug("setting up watchers func")
	return func(context context.Context) {
		l.Debug("watching for resources")
		list, err := c.List(context, v1.ListOptions{})
		s.Set(list)
		if err != nil {
			l.Errorf("when watching for shoots: %s", err.Error())
			return
		}

		l.Debugf("found %v shoots", len(list.Items))
	}
}
