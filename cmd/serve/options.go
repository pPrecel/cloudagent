package serve

import (
	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pkg/errors"
)

type options struct {
	*command.Options

	configPath string

	onDemand      bool
	socketAddress string
	socketNetwork string
}

func NewOptions(opts *command.Options) *options {
	return &options{
		Options:       opts,
		socketNetwork: agent.Network,
	}
}

func (o *options) validate() error {
	if o.configPath == "" {
		return errors.New("configPath should not be empty")
	}

	return nil
}
