package main

import (
	"context"
	"os"

	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/cmd/generate"
	"github.com/pPrecel/cloudagent/cmd/serve"
	"github.com/pPrecel/cloudagent/cmd/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	o := &command.Options{
		Context: context.Background(),
		Logger:  newLogger(),
	}

	cmd := &cobra.Command{
		Use:          "cloudagent",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if o.Verbose == true {
				o.Logger.SetLevel(logrus.DebugLevel)
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(&o.Verbose, "verbose", "v", false, "Displays details of actions triggered by the command.")

	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	cmd.AddCommand(&cobra.Command{Use: "completion", Hidden: true})

	cmd.AddCommand(
		generate.NewCmd(o),
		serve.NewCmd(serve.NewOptions(o)),
		state.NewCmd(state.NewOptions(o)),
	)

	err := cmd.Execute()
	if err != nil {
		o.Logger.Fatal(err)
	}
}

func newLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)

	return log
}
