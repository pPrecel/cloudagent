package generate

import (
	"errors"
	"io"
	"os"

	command "github.com/pPrecel/cloud-agent/cmd"
)

type options struct {
	*command.Options

	configPath   string
	agentVerbose bool

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
	if o.configPath == "" {
		return errors.New("configPath should not be empty")
	}

	return nil
}
