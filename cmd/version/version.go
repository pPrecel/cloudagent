package version

import (
	"fmt"

	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/spf13/cobra"
)

func NewCmd(o *command.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf(o.Version)
		},
	}

	return cmd
}
