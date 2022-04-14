package serve

import (
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pPrecel/gardener-agent/internal/gardener"
	"github.com/spf13/cobra"
)

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "serve",
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(o)
		},
	}

	return cmd
}

func run(o *options) error {
	o.Logger.Info("starting gardeners agent daemon")

	o.Logger.Info("loading configuration...")
	// TODO: load config

	o.Logger.Info("creating gardeners client")
	cfg, err := gardener.NewClusterConfig("/Users/i517616/Desktop/kubeconfig-garden-wookiee.yml")
	if err != nil {
		o.Logger.Fatal(err)
	}

	client, err := gardener.NewClient(cfg)
	if err != nil {
		o.Logger.Fatal(err)
	}

	list, err := client.Shoots("garden-wookiee").List(o.Ctx, v1.ListOptions{})
	if err != nil {
		o.Logger.Fatal(err)
	}

	for i := range list.Items {
		item := list.Items[i]
		fmt.Printf("\n\n%v: %s,\n %+v\n\n", i, item.Name, item.ObjectMeta.Annotations)
	}

	o.Logger.Info("starting state watcher")
	// TODO: start cron job

	o.Logger.Info("starting grpc socket server")
	// TODO: start server
}
