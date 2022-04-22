package state

import (
	"time"

	"github.com/pPrecel/cloud-agent/internal/command"
	"github.com/pPrecel/cloud-agent/internal/output"
)

type options struct {
	*command.Options
	CreatedBy       string
	StringOutFormat string
	StringErrFormat string
	OutFormat       output.Output
	Timeout         time.Duration
}

func NewOptions(opts *command.Options) *options {
	return &options{
		Options: opts,
	}
}
