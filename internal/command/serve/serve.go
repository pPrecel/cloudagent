package serve

import (
	"github.com/pPrecel/cloud-agent/internal/gardener"
	"github.com/pPrecel/cloud-agent/pkg/agent"
	cloud_agent "github.com/pPrecel/cloud-agent/pkg/agent/proto"
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

	cmd.Flags().StringVarP(&o.kubeconfigPath, "kubeconfigPath", "k", "", "Provides path to kubeconfig.")
	cmd.Flags().StringVarP(&o.namespace, "namespace", "n", "", "Provides gardener namespace.")
	cmd.Flags().StringVarP(&o.cronSpec, "cronSpec", "c", "@every 15m", "Provices spec for cron configuration.")

	return cmd
}

func run(o *options) error {
	o.Logger.Info("starting gardeners agent")

	o.Logger.Debug("creating cluster config")
	cfg, err := o.newClusterConfig(o.kubeconfigPath)
	if err != nil {
		return err
	}

	o.Logger.Debug("creating gardener client")
	c, err := gardener.NewClient(cfg)
	if err != nil {
		return err
	}

	state := &gardener.LastState{}

	o.Logger.Infof("starting state watcher with spec: '%s'", o.cronSpec)
	gardenerFn := o.newWatchFunc(o.Logger, c.Shoots(o.namespace), state)

	watcher, err := agent.NewWatcher(agent.WatcherOptions{
		Spec:    o.cronSpec,
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
	lis, err := agent.NewSocket(o.socketNetwork, o.socketAddress)
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
