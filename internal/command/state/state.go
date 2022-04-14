package state

import "github.com/spf13/cobra"

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "state",
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(o)
		},
	}

	return cmd
}

func run(o *options) error {
	return nil
}
