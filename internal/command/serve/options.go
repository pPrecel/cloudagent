package serve

import "github.com/pPrecel/gardener-agent/internal/command"

type options struct {
	*command.Options
	KubeconfigPath string
	Namespace      string
}

func NewOptions(opts *command.Options) *options {
	return &options{
		Options: opts,
	}
}
