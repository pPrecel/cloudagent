package logs

import (
	"github.com/hpcloud/tail"
	"github.com/pPrecel/cloudagent/pkg/brew"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get cloudagent logs.",
		Long:  "Use this command to print or follow cloudagent logs.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(o)
		},
	}

	cmd.Flags().StringVar(&o.filePath, "file", brew.StdoutPath, "Provides path to the 'cloudagent.stdout' file.")
	cmd.Flags().BoolVarP(&o.followLogs, "follow", "f", false, "Follow logs.")

	return cmd
}

func run(o *options) error {
	t, err := o.tailFile(o.filePath, tail.Config{
		Follow: o.followLogs,
	})
	if err != nil {
		return errors.Wrapf(err, "can't tail '%s' file", o.filePath)
	}

	for l := range t.Lines {
		_, err = o.writer.Write(append([]byte(l.Text), '\n'))
		if err != nil {
			return errors.Wrapf(err, "can't print line: %s", l.Text)
		}
	}

	return nil
}
