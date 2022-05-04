package schema

import (
	"io"
	"os"

	command "github.com/pPrecel/cloud-agent/cmd"
	"github.com/pPrecel/cloud-agent/pkg/config"
)

type options struct {
	*command.Options

	stdout     io.Writer
	jsonSchema func() ([]byte, error)
}

func NewOptions(o *command.Options) *options {
	return &options{
		Options:    o,
		stdout:     os.Stdout,
		jsonSchema: config.JSONSchema,
	}
}
