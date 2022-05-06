package serve

import (
	"net"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/internal/gardener"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

type options struct {
	*command.Options

	configPath string

	socketAddress    string
	socketNetwork    string
	getConfig        func(string) (*config.Config, error)
	newClusterConfig func(string) (*rest.Config, error)
	newSocket        func(network, address string) (net.Listener, error)
	newWatchFunc     func(l *logrus.Logger, c gardener.Client, s agent.RegisteredResource[*v1beta1.ShootList]) agent.WatchFn
}

func NewOptions(opts *command.Options) *options {
	return &options{
		Options:          opts,
		socketAddress:    agent.Address,
		socketNetwork:    agent.Network,
		getConfig:        config.GetConfig,
		newClusterConfig: gardener.NewClusterConfig,
		newSocket:        agent.NewSocket,
		newWatchFunc:     gardener.NewWatchFunc,
	}
}

func (o *options) validate() error {
	if o.configPath == "" {
		return errors.New("configPath should not be empty")
	}

	return nil
}
