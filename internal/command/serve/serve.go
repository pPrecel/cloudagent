package serve

import (
	"github.com/pPrecel/cloud-agent/internal/agent"
	cloud_agent "github.com/pPrecel/cloud-agent/internal/agent/proto"
	"github.com/pPrecel/cloud-agent/internal/gardener"
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

	state := &gardener.LastState{}

	o.Logger.Infof("starting state watcher with spec: '%s'", o.CronSpec)
	gardenerFn, err := gardener.NewWatchFunc(gardener.WatchOptions{
		KubeconfigPath: o.KubeconfigPath,
		Namespace:      o.Namespace,
		StateSetter:    state,
		Logger:         o.Logger,
	})
	if err != nil {
		return err
	}

	watcher, err := agent.NewWatcher(agent.WatcherOptions{
		Spec:    o.CronSpec,
		Context: o.Context,
		Logger:  o.Logger,
	}, gardenerFn)
	if err != nil {
		return err
	}
	defer watcher.Stop()

	o.Logger.Debug("starting watcher")
	watcher.Start()

	o.Logger.Debug("configuring grpc server")
	lis, err := agent.NewSocket(agent.Network, agent.Address)
	if err != nil {
		return err
	}

	grpcServer := googlerpc.NewServer(googlerpc.EmptyServerOption{})
	agentServer := agent.NewServer(&agent.ServerOption{
		Getter: state,
		Logger: o.Logger,
	})
	cloud_agent.RegisterAgentServer(grpcServer, agentServer)

	o.Logger.Info("starting grpc server")
	return grpcServer.Serve(lis)
}
