package main

import (
	"context"
	"os"

	"github.com/pPrecel/gardener-agent/internal/command"
	"github.com/pPrecel/gardener-agent/internal/command/generate"
	"github.com/pPrecel/gardener-agent/internal/command/serve"
	"github.com/pPrecel/gardener-agent/internal/command/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)

	o := &command.Options{
		Context: context.Background(),
		Logger:  log,
	}

	cmd := &cobra.Command{
		Use:          "gardenagent",
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
		serve.NewCmd(serve.NewOptions(o)),
		state.NewCmd(state.NewOptions(o)),
		generate.NewCmd(generate.NewOptions(o)),
	)

	err := cmd.Execute()
	if err != nil {
		o.Logger.Fatal(err)
	}
}
