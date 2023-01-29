package serve

import (
	"github.com/pPrecel/cloudagent/internal/watcher"
	"github.com/pPrecel/cloudagent/pkg/agent"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	googlerpc "google.golang.org/grpc"
)

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve clouds watching.",
		Long:  "Use this command to serve an agent functionality to observe clouds you specify in the configuration file.",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return o.validate()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			// change logger options
			o.Logger.Formatter = &logrus.JSONFormatter{}

			return run(o)
		},
	}

	cmd.Flags().StringVarP(&o.configPath, "config-path", "c", config.ConfigPath, "Provides path to the config file.")
	cmd.Flags().StringVar(&o.socketAddress, "socket-path", agent.Address, "Provides path to the socket file.")

	return cmd
}

func run(o *options) error {
	o.Logger.Info("starting gardeners agent")

	resourceGetter, err := watcher.NewForConfig(&watcher.Options{
		Context:    o.Context,
		Logger:     o.Logger.WithField("component", "watcher"),
		ConfigPath: o.configPath,
	})
	if err != nil {
		return err
	}

	o.Logger.Debugf("configuring grpc server - network '%s', address '%s'", o.socketNetwork, o.socketAddress)
	lis, err := agent.NewSocket(o.socketNetwork, o.socketAddress)
	if err != nil {
		return err
	}

	grpcServer := googlerpc.NewServer(googlerpc.EmptyServerOption{})
	agentServer := agent.NewServer(&agent.ServerOption{
		ResourceGetter: resourceGetter,
		Logger:         o.Logger.WithField("component", "server"),
	})
	cloud_agent.RegisterAgentServer(grpcServer, agentServer)

	o.Logger.Info("starting grpc server")
	return grpcServer.Serve(lis)
}
