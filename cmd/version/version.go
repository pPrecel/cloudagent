package version

import (
	"fmt"

	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/spf13/cobra"
)

func NewCmd(o *command.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		RunE: func(_ *cobra.Command, _ []string) error {
			_, err := fmt.Printf(o.Version)
			return err
		},
	}

	return cmd
}
