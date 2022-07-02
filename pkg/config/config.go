package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	ConfigFilename = ".cloudagent.conf.yaml"
	fileMode       = 0644
)

var (
	ConfigPath = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ConfigFilename)
)

func Read(path string) (*Config, error) {
	cfg := Config{}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read config from file: '%s'", path)
	}

	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal config")
	}

	return &cfg, nil
}

func Write(path string, config interface{}) error {
	return write(yaml.Marshal, path, config)
}

func write(marshal func(interface{}) ([]byte, error), path string, config interface{}) error {
	b, err := marshal(config)
	if err != nil {
		return errors.Wrapf(err, "can't write config to file: '%s'", path)
	}

	return ioutil.WriteFile(path, b, fileMode)
}
