package main

import (
	"context"
	"os"

	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/cmd/check"
	"github.com/pPrecel/cloudagent/cmd/config"
	"github.com/pPrecel/cloudagent/cmd/logs"
	"github.com/pPrecel/cloudagent/cmd/serve"
	"github.com/pPrecel/cloudagent/cmd/state"
	"github.com/pPrecel/cloudagent/cmd/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Version = "local"

func main() {
	o := &command.Options{
		Version: Version,
		Context: context.Background(),
		Logger:  newLogger(),
	}

	cmd := &cobra.Command{
		Use:          "cloudagent",
		SilenceUsage: true,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			if o.Verbose {
				o.Logger.SetLevel(logrus.DebugLevel)
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(&o.Verbose, "verbose", "v", false, "Displays details of actions triggered by the command.")

	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	cmd.AddCommand(&cobra.Command{Use: "completion", Hidden: true})

	cmd.AddCommand(
		config.NewCmd(config.NewOptions(o)),
		serve.NewCmd(serve.NewOptions(o)),
		state.NewCmd(state.NewOptions(o)),
		check.NewCmd(check.NewOptions(o)),
		logs.NewCmd(logs.NewOptions(o)),
		version.NewCmd(o),
	)

	_ = cmd.Execute()
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
