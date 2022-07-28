package state

import (
	"fmt"
	"time"

	"github.com/pPrecel/cloudagent/internal/formater"
	"github.com/pPrecel/cloudagent/internal/output"
	"github.com/pPrecel/cloudagent/pkg/agent"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	warningLog = "Warning: Some data may be not up to date because of cloudagent error.\nTry 'cloudagent check' for more info.\n\n"
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
- '`+formater.GardenerTextHealthyFormat+`' represents number of clusters with the HEALTHY status,
- '`+formater.GardenerTextHibernatedFormat+`' represents number of clusters with the HIBERNATED status,
- '`+formater.GardenerTextUnknownFormat+`' represents number of clusters with the UNKNOWN status,
- '`+formater.GardenerTextEmptyFormat+`' represents number of clusters with the EMPTY status,
- '`+formater.GardenerTextEmptyUnknownFormat+`' represents number of clusters with the EMPTY or the UNKNOWN status,
- '`+formater.GardenerTextAllFormat+`' represents of all clusters in namespace.

The second one can contains '`+formater.GardenerTextErrorFormat+`'  which will be replaced with error message.`)
	cmd.Flags().DurationVarP(&o.timeout, "timeout", "t", 2*time.Second, "Provides timeout for the command.")
	cmd.Flags().StringVar(&o.socketAddress, "socket-path", agent.Address, "Provides path to the socket file.")

	return cmd
}

func run(o *options) error {
	o.Logger.Debugf("getting shoots")
	list, err := shootState(o)

	o.Logger.Debugf("received: %+v, error: %v", list, err)

	if err != nil {
		return errors.Wrap(err, "cloudagent internal error")
	}

	f := formater.NewGardener(list.ShootList, formater.Filters{
		CreatedBy: o.createdBy,
	})

	// print warning
	if isAnyError(o.Logger, list) &&
		o.outFormat.Type() == string(output.TableType) {
		o.Logger.Debug("printing warning log")
		o.writer.Write([]byte(warningLog))
	}

	return o.outFormat.Print(o.writer, f)
}

func shootState(o *options) (*cloud_agent.GardenerResponse, error) {
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

func isAnyError(l *logrus.Logger, resp *cloud_agent.GardenerResponse) bool {
	if resp.GeneralError != "" {
		l.Debugf("got general error: '%s'", resp.GeneralError)
		return true
	}

	for key := range resp.ShootList {
		shootList := resp.ShootList[key]
		if shootList != nil && shootList.Error != "" {
			l.Debugf("got error: '%s'", shootList.Error)
			return true
		}
	}

	return false
}
