package config

import (
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manipulate cloudagent configuration.",
		Long:  "Use this command to extend cloudagent configuration file in a more user-friendly way.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(o)
		},
	}

	cmd.AddCommand(
		newSchemaCmd(newSchemaOptions(o)),
		newGardenerCmd(newGardenerOptions(o)),
	)

	cmd.PersistentFlags().StringVarP(&o.configPath, "config-path", "c", config.ConfigPath, "Provides path to the config file.")

	return cmd
}

func run(o *options) error {
	c, err := o.readConfig(o.configPath)
	if err != nil {
		return errors.Wrap(err, "while reading config file")
	}

	b, err := yaml.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "while verifying config structure")
	}

	_, err = o.stdout.Write(b)
	if err != nil {
		return errors.Wrap(err, "while printing config file")
	}

	return nil
}
