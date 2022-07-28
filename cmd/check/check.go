package check

import (
	"context"
	"fmt"
	"time"

	"github.com/pPrecel/cloudagent/internal/formater"
	"github.com/pPrecel/cloudagent/internal/output"
	"github.com/pPrecel/cloudagent/pkg/agent"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check status of the server.",
		Long:  "Use this command to check if the server faces any problem with config, credentials, clouds, and so on.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(o)
		},
	}

	cmd.Flags().VarP(output.NewFlag(&o.outFormat, "table", "$g/$G/$a", "-/-/-/-"), "output", "o", ``)
	cmd.Flags().DurationVarP(&o.timeout, "timeout", "t", 2*time.Second, "Provides timeout for the command.")
	cmd.Flags().StringVar(&o.socketAddress, "socket-path", agent.Address, "Provides path to the socket file.")

	return cmd
}

func run(o *options) error {
	o.Logger.Debugf("getting shoots")
	list, err := shootState(o)
	if err != nil {
		return err
	}

	// print warning
	if list != nil && list.GeneralError != "" &&
		o.outFormat.ErrorFormat() == string(output.TableType) {
		o.Logger.Debug("printing warning log")
		o.stdout.Write([]byte(list.GeneralError))
	}

	f := formater.NewCheck(list)
	return o.outFormat.Print(o.stdout, f)
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
	resp, err := cloud_agent.NewAgentClient(conn).GardenerShoots(ctx, &cloud_agent.Empty{})
	if err != nil {
		o.Logger.Debugf("fail to get shoots: %v", err)
		return nil, err
	}

	return resp, nil
}
