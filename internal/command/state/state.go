package state

import (
	"fmt"
	"time"

	"github.com/pPrecel/cloud-agent/internal/agent"
	cloud_agent "github.com/pPrecel/cloud-agent/internal/agent/proto"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "state",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return o.validate()
		},
		Run: func(_ *cobra.Command, _ []string) {
			run(o)
		},
	}

	cmd.Flags().StringVarP(&o.CreatedBy, "createdBy", "c", "", "Provides filter argument for owned, hibernated and corrupted shoots.")
	cmd.Flags().StringVarP(&o.OutFormat, "output-format", "o", "%d/%d/%d/%d", `Provides format for the output information. Must contains four '%d' elements where:
	- first is number of running clusters, 
	- second is number of hibernated clusters, 
	- third is number of cluster with unknown status, 
	- fourth is number of all cluster in namespace.`)
	cmd.Flags().StringVarP(&o.ErrFormat, "error-format", "e", "ERR", "Provides format of output after occures an error.")
	cmd.Flags().DurationVarP(&o.Timeout, "timeout", "t", 2*time.Second, "Provides timeout for the command.")

	return cmd
}

func run(o *options) {
	list, err := shootState(o)
	if err != nil {
		fmt.Print(o.ErrFormat)
		return
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

		if list.Shoots[i].Condition == cloud_agent.Condition_HIBERNATED {
			hibernated++
		} else if list.Shoots[i].Condition == cloud_agent.Condition_UNKNOWN {
			corrupted++
		} else {
			healthy++
		}
	}

	fmt.Printf(o.OutFormat, healthy, hibernated, corrupted, len(list.Shoots))
}

func shootState(o *options) (*cloud_agent.ShootList, error) {
	o.Logger.Debug("creating grpc client")
	conn, err := grpc.Dial(fmt.Sprintf("%s://%s", agent.Network, agent.Address), grpc.WithInsecure())
	if err != nil {
		o.Logger.Debugf("fail to dial: %v", err)
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(o.Context, o.Timeout)
	defer cancel()

	o.Logger.Debug("sending request")
	list, err := cloud_agent.NewAgentClient(conn).GardenerShoots(ctx, &cloud_agent.Empty{})
	if err != nil {
		o.Logger.Debugf("fail to get shoots: %v", err)
		return nil, err
	}

	return list, nil
}

func isCreatedBy(creator string, shoot *cloud_agent.Shoot) bool {
	if shoot.Annotations["gardener.cloud/created-by"] == creator {
		return true
	}
	return false
}
