package state

import "github.com/pPrecel/gardener-agent/internal/command"

type options struct {
	*command.Options
	selector string
	format   string
}

func NewOptions(opts *command.Options) *options {
	return &options{
		Options: opts,
	}
}
