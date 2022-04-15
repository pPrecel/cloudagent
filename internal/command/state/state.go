package state

import (
	"fmt"

	"github.com/pPrecel/gardener-agent/internal/agent"
	gardener_agent "github.com/pPrecel/gardener-agent/internal/agent/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "state",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return o.validate()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(o)
		},
	}

	cmd.Flags().StringVarP(&o.CreatedBy, "createdBy", "c", "", "Provides filter argument for owned, hibernated and corrupted shoots.")

	return cmd
}

func run(o *options) error {
	o.Logger.Debug("creating grpc client")
	conn, err := grpc.Dial(fmt.Sprintf("%s://%s", agent.Network, agent.Address), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logrus.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	o.Logger.Debug("sending request")
	list, err := gardener_agent.NewAgentClient(conn).Shoots(o.Context, &gardener_agent.Empty{})
	if err != nil {
		return err
	}

	hibernated := 0
	corrupted := 0
	healthy := 0
	o.Logger.Debug("received %v items", len(list.Shoots))
	for i := range list.Shoots {
		o.Logger.Debug("%v - %+v", i, list.Shoots[i])
		if !isCreatedBy(o.CreatedBy, list.Shoots[i]) {
			continue
		}

		if list.Shoots[i].Condition == gardener_agent.Condition_HIBERNATED {
			hibernated++
		} else if list.Shoots[i].Condition == gardener_agent.Condition_UNKNOWN {
			corrupted++
		} else {
			healthy++
		}
	}

	fmt.Printf("%d/%d/%d/%d", healthy, hibernated, corrupted, len(list.Shoots))

	return nil
}

func isCreatedBy(creator string, shoot *gardener_agent.Shoot) bool {
	if shoot.Annotations["gardener.cloud/created-by"] == creator {
		return true
	}
	return false
}
