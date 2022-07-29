package config

import (
	"github.com/pPrecel/cloudagent/internal/formater"
	"github.com/pPrecel/cloudagent/internal/output"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manipulate cloudagent configuration.",
		Long:  "Use this command to extend cloudagent configuration file in a more user-friendly way.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(o)
		},
		Example: `  # Get config
  cloudagent config

  # Get config with custom test output
  cloudagent config -o text=$a=$E`,
	}

	cmd.AddCommand(
		newSchemaCmd(newSchemaOptions(o)),
		newGardenerCmd(newGardenerOptions(o)),
	)

	cmd.PersistentFlags().StringVarP(&o.configPath, "config-path", "c", config.ConfigPath, "Provides path to the config file.")
	cmd.Flags().VarP(output.NewFlag(&o.outFormat, "table", "$g/$G/$a", "-/-/-/-"), "output", "o", `Provides format for the output information. 
	
For the 'text' output format you can specifie two more informations by spliting them using '='. The first one would be used as output format and second as error format.

The first one can contains at least on out of four elements where:
- '`+formater.ConfigTextAllFormat+`' represents number of all projects,
- '`+formater.ConfigTextGCPFormat+`' represents number of GCP projects,
- '`+formater.ConfigTextGardenerFormat+`' represents number of Gardener projects,
- '`+formater.ConfigTextPersistentFormat+`' represents value of the persistentSpec field.

The second one can contains '`+formater.ConfigTextErrorFormat+`'  which will be replaced with error message.`)

	return cmd
}

func run(o *options) error {
	c, err := o.readConfig(o.configPath)
	if err != nil {
		return errors.Wrap(err, "while reading config file")
	}

	// verify
	if _, err := o.marshal(c); err != nil {
		return errors.Wrap(err, "while verifying config structure")
	}

	return o.outFormat.Print(o.stdout, formater.NewConfig(err, c))
}
