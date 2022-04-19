package serve

import (
	"path/filepath"

	"github.com/pPrecel/cloud-agent/internal/command"
	"github.com/pkg/errors"
)

type options struct {
	*command.Options
	KubeconfigPath string
	Namespace      string
	CronSpec       string
}

func NewOptions(opts *command.Options) *options {
	return &options{
		Options: opts,
	}
}

func (o *options) validate() error {
	if !filepath.IsAbs(o.KubeconfigPath) {
		path, err := filepath.Abs(o.KubeconfigPath)
		if err != nil {
			return errors.Wrap(err, "kubeconfigPath should not be empty")
		}

		o.KubeconfigPath = path
	}

	if o.Namespace == "" {
		return errors.New("namespace should not be empty")
	}

	return nil
}
