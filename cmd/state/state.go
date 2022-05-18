package state

import (
	"fmt"
	"time"

	"github.com/pPrecel/cloudagent/internal/formater"
	"github.com/pPrecel/cloudagent/internal/output"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Get cached states from the agent.",
		Long:  "Use this command to communicate with the agent's socket and take the info about clouds in the specified output type.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(o)
		},
	}

	cmd.Flags().StringVarP(&o.createdBy, "created-by", "c", "", "Show clusters created by specific person.")
	cmd.Flags().VarP(output.NewFlag(&o.outFormat, "table", "$r/$h/$x/$a", "-/-/-/-"), "output", "o", `Provides format for the output information. 
	
For the 'text' output format you can specifie two more informations by spliting them using '='. The first one would be used as output format and second as error format.

The first one can contains at least on out of four elements where:
- '`+formater.TextHealthyFormat+`' represents number of clusters with the HEALTHY status,
- '`+formater.TextHibernatedFormat+`' represents number of clusters with the HIBERNATED status,
- '`+formater.TextUnknownFormat+`' represents number of clusters with the UNKNOWN status,
- '`+formater.TextEmptyFormat+`' represents number of clusters with the EMPTY status,
- '`+formater.TextEmptyUnknownFormat+`' represents number of clusters with the EMPTY or the UNKNOWN status,
- '`+formater.TextAllFormat+`' represents of all clusters in namespace.

The second one can contains '`+formater.TextErrorFormat+`'  which will be replaced with error message.`)
	cmd.Flags().DurationVarP(&o.timeout, "timeout", "t", 2*time.Second, "Provides timeout for the command.")

	return cmd
}

func run(o *options) error {
	o.Logger.Debugf("getting shoots")
	list, err := shootState(o)

	o.Logger.Debugf("received: %+v, error: %v", list, err)

	f := formater.NewForState(err, list, formater.Filters{
		CreatedBy: o.createdBy,
	})

	return o.outFormat.Print(o.writer, f)
}

func shootState(o *options) (*cloud_agent.ShootList, error) {
	o.Logger.Debug("creating grpc client")
	conn, err := grpc.Dial(fmt.Sprintf("%s://%s", o.socketNetwork, o.socketAddress), grpc.WithInsecure())
	if err != nil {
		o.Logger.Debugf("fail to dial: %v", err)
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(o.Context, o.timeout)
	defer cancel()

	o.Logger.Debug("sending request")
	list, err := cloud_agent.NewAgentClient(conn).GardenerShoots(ctx, &cloud_agent.Empty{})
	if err != nil {
		o.Logger.Debugf("fail to get shoots: %v", err)
		return nil, err
	}

	return list, nil
}
