package state

import (
	"time"

	"github.com/pPrecel/gardener-agent/internal/command"
	"github.com/pkg/errors"
)

type options struct {
	*command.Options
	CreatedBy string
	OutFormat string
	ErrFormat string
	Timeout   time.Duration
}

func NewOptions(opts *command.Options) *options {
	return &options{
		Options: opts,
	}
}

func (o *options) validate() error {
	if o.CreatedBy == "" {
		return errors.New("createdBy should not be empty")
	}

	return nil
}