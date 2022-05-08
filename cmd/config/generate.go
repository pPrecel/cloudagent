package config

import (
	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/cmd/config/schema"
	"github.com/spf13/cobra"
)

func NewCmd(o *command.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manipulate cloudagent configuration.",
		Long:  "Use this command to extend cloudagent configuration file in a more user-friendly way.",
	}

	cmd.AddCommand(
		schema.NewCmd(schema.NewOptions(o)),
	)

	return cmd
}
