package config

import (
	"io"
	"os"

	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/internal/output"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type options struct {
	*command.Options

	outFormat  output.Flag
	configPath string

	stdout      io.Writer
	marshal     func(interface{}) ([]byte, error)
	readConfig  func(string) (*config.Config, error)
	writeConfig func(string, interface{}) error
}

func NewOptions(o *command.Options) *options {
	return &options{
		Options:     o,
		stdout:      os.Stdout,
		marshal:     yaml.Marshal,
		readConfig:  config.Read,
		writeConfig: config.Write,
	}
}

type schemaOptions struct {
	*options

	stdout     io.Writer
	jsonSchema func() ([]byte, error)
}

func newSchemaOptions(o *options) *schemaOptions {
	return &schemaOptions{
		options:    o,
		stdout:     os.Stdout,
		jsonSchema: config.JSONSchema,
	}
}

type gardenerOptions struct {
	*options

	namespace  string
	kubeconfig string
}

func newGardenerOptions(o *options) *gardenerOptions {
	return &gardenerOptions{
		options: o,
	}
}

func (o *gardenerOptions) validateAdd() error {
	if o.namespace == "" || o.kubeconfig == "" {
		return errors.New("namespace and kubeconfig can't be empty")
	}

	return nil
}

func (o *gardenerOptions) validateDel() error {
	if o.namespace == "" && o.kubeconfig == "" {
		return errors.New("namespace or kubeconfig can't be empty")
	}

	return nil
}
