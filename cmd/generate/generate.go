package generate

import (
	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/cmd/generate/plist"
	"github.com/pPrecel/cloudagent/cmd/generate/schema"
	"github.com/spf13/cobra"
)

func NewCmd(o *command.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate some system-oriented utils.",
		Long:  "Use this command to generate a system-oriented utils it would help you in communication between the cloudagent and other tools.",
	}

	cmd.AddCommand(
		plist.NewCmd(plist.NewOptions(o)),
		schema.NewCmd(schema.NewOptions(o)),
	)

	return cmd
}
