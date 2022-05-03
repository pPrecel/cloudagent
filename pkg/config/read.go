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
)

var (
	ConfigPath = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ConfigFilename)
)

func GetConfig(path string) (*Config, error) {
	cfg := Config{}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read config file: '%s'", path)
	}

	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal config")
	}

	return &cfg, nil
}
