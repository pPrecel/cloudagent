package serve

import (
	v1beta1_apis "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloud-agent/internal/gardener"
	"github.com/pPrecel/cloud-agent/pkg/agent"
	cloud_agent "github.com/pPrecel/cloud-agent/pkg/agent/proto"
	"github.com/pPrecel/cloud-agent/pkg/config"
	"github.com/robfig/cron/v3"
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

	cmd.Flags().StringVarP(&o.configPath, "configPath", "c", config.ConfigPath, "Provides path to the config file.")

	return cmd
}

func run(o *options) error {
	o.Logger.Info("starting gardeners agent")

	o.Logger.Infof("reading config from path: '%s'", o.configPath)
	cfg, err := o.getConfig(o.configPath)
	if err != nil {
		return err
	}

	gardenerCache := agent.NewCache[*v1beta1_apis.ShootList]()

	o.Logger.Infof("starting state watcher with spec: '%s'", cfg.PersistentSpec)
	watcher, err := buildWatcher(o, cfg, gardenerCache)
	if err != nil {
		return err
	}
	defer watcher.Stop()

	o.Logger.Debug("starting watcher")
	watcher.Start()

	o.Logger.Debug("configuring grpc server")
	lis, err := o.newSocket(o.socketNetwork, o.socketAddress)
	if err != nil {
		return err
	}

	grpcServer := googlerpc.NewServer(googlerpc.EmptyServerOption{})
	agentServer := agent.NewServer(&agent.ServerOption{
		GardenerCache: gardenerCache,
		Logger:        o.Logger,
	})
	cloud_agent.RegisterAgentServer(grpcServer, agentServer)

	o.Logger.Info("starting grpc server")
	return grpcServer.Serve(lis)
}

func buildWatcher(o *options, config *config.Config, cache agent.Cache[*v1beta1_apis.ShootList]) (*cron.Cron, error) {
	funcs := []agent.WatchFn{}
	for i := range config.GardenerProjects {
		p := config.GardenerProjects[i]

		o.Logger.Debug("creating cluster config")
		cfg, err := o.newClusterConfig(p.KubeconfigPath)
		if err != nil {
			return nil, err
		}

		o.Logger.Debug("creating gardener client")
		c, err := gardener.NewClient(cfg)
		if err != nil {
			return nil, err
		}

		r := cache.Register(config.GardenerProjects[i].Namespace)

		o.Logger.Debugf("creeating watcher func for namespace: '%s'", p.Namespace)
		funcs = append(funcs,
			o.newWatchFunc(o.Logger, c.Shoots(config.GardenerProjects[i].Namespace), r),
		)
	}

	return agent.NewWatcher(agent.WatcherOptions{
		Spec:    config.PersistentSpec,
		Context: o.Context,
		Logger:  o.Logger,
	}, funcs...)
}
