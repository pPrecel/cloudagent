package generate

import (
	"errors"
	"io"
	"os"

	"github.com/pPrecel/cloud-agent/internal/command"
)

type options struct {
	*command.Options

	kubeconfigPath string
	namespace      string
	cronSpec       string
	agentVerbose   bool

	executable func() (string, error)
	stdout     io.Writer
}

func NewOptions(o *command.Options) *options {
	return &options{
		Options:    o,
		executable: os.Executable,
		stdout:     os.Stdout,
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
