package main

import (
	"context"

	"github.com/pPrecel/gardener-agent/internal/command"
	"github.com/pPrecel/gardener-agent/internal/command/serve"
	"github.com/pPrecel/gardener-agent/internal/command/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "gardenagent",
	}

	o := &command.Options{
		Ctx:    context.Background(),
		Logger: logrus.New(),
	}

	cmd.PersistentFlags().BoolVarP(&o.Verbose, "verbose", "v", false, "Displays details of actions triggered by the command.")
	cmd.PersistentFlags().BoolP("help", "h", false, "Provides command help.")

	cmd.AddCommand(
		serve.NewCmd(serve.NewOptions(o)),
		state.NewCmd(state.NewOptions(o)),
	)

	err := cmd.Execute()
	if err != nil {

	}
}
