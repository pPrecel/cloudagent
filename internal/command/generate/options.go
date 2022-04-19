package generate

import "github.com/pPrecel/cloud-agent/internal/command"

type options struct {
	*command.Options
	KubeconfigPath string
	Namespace      string
	CronSpec       string
	AgentVerbose   bool
}

func NewOptions(o *command.Options) *options {
	return &options{
		Options: o,
	}
}
