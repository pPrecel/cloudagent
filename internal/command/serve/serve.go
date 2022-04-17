package serve

import (
	"net"
	"os"

	"github.com/pPrecel/gardener-agent/internal/agent"
	gardener_agent "github.com/pPrecel/gardener-agent/internal/agent/proto"
	"github.com/spf13/cobra"
	googlerpc "google.golang.org/grpc"
)

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "serve",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return o.validate()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(o)
		},
	}

	cmd.Flags().StringVarP(&o.KubeconfigPath, "kubeconfigPath", "k", "", "Provides path to kubeconfig.")
	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", "", "Provides gardener namespace.")
	cmd.Flags().StringVarP(&o.CronSpec, "cronSpec", "c", "@every 15m", "Provices spec for cron configuration.")

	return cmd
}

func run(o *options) error {
	o.Logger.Info("starting gardeners agent")

	o.Logger.Infof("removing old socket: '%s'", agent.Address)
	err := os.RemoveAll(agent.Address)
	if err != nil {
		return err
	}

	state := &agent.LastState{}

	o.Logger.Infof("starting state watcher with spec: '%s'", o.CronSpec)
	watcher, err := agent.NewWatcher(agent.WatcherOption{
		KubeconfigPath: o.KubeconfigPath,
		Namespace:      o.Namespace,
		Spec:           o.CronSpec,
		Context:        o.Context,
		StateSetter:    state,
		Logger:         o.Logger,
	})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	o.Logger.Debug("starting watcher")
	watcher.Start()

	o.Logger.Debug("configuring grpc server")
	lis, err := net.Listen(agent.Network, agent.Address)
	if err != nil {
		return err
	}

	grpcServer := googlerpc.NewServer(googlerpc.EmptyServerOption{})
	agentServer := agent.NewServer(&agent.ServerOption{
		Getter: state,
		Logger: o.Logger,
	})
	gardener_agent.RegisterAgentServer(grpcServer, agentServer)

	o.Logger.Info("starting grpc server")
	return grpcServer.Serve(lis)
}
