package serve

import (
	"github.com/pPrecel/cloud-agent/internal/command"
	"github.com/pPrecel/cloud-agent/internal/gardener"
	"github.com/pPrecel/cloud-agent/pkg/agent"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

type options struct {
	*command.Options

	kubeconfigPath string
	namespace      string
	cronSpec       string

	socketAddress    string
	socketNetwork    string
	newClusterConfig func(string) (*rest.Config, error)
	newWatchFunc     func(l *logrus.Logger, c gardener.Client, s gardener.StateSetter) agent.WatchFn
}

func NewOptions(opts *command.Options) *options {
	return &options{
		Options:          opts,
		socketAddress:    agent.Address,
		socketNetwork:    agent.Network,
		newClusterConfig: gardener.NewClusterConfig,
		newWatchFunc:     gardener.NewWatchFunc,
	}
}

func (o *options) validate() error {
	if o.kubeconfigPath == "" {
		return errors.New("kubeconfigPath should not be empty")
	}

	if o.namespace == "" {
		return errors.New("namespace should not be empty")
	}

	return nil
}
