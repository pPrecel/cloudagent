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
		Example: `  # Check agent
  cloudagent check
  
  # Check with text output
  cloudagent check -o text
  
  # Check with custom test output
  cloudagent check -o text=$a=$E`,
	}

	cmd.Flags().VarP(output.NewFlag(&o.outFormat, "table", "$h/$e/$a", "$E"), "output", "o", `Provides format for the output information. 
	
	For the 'text' output format you can specifie two more informations by spliting them using '='. The first one would be used as output format and second as error format.
	
	The first one can contains at least on out of four elements where:
	- '`+formater.CheckTextAllFormat+`' represents number of all projects,
	- '`+formater.CheckTextHealthyFormat+`' represents number of all healthy projects,
	- '`+formater.CheckTextErrorCountFormat+`' represents number of all projects with error,
	- '`+formater.CheckTextErrorFormat+`' represents error message.

	
	The second one can contains '`+formater.CheckTextErrorFormat+`'  which will be replaced with error message when response is nil.`)
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
		o.outFormat.Type() == string(output.TableType) {
		o.Logger.Debug("printing warning log")
		o.stdout.Write([]byte(list.GeneralError))
	}

	f := formater.NewCheck(list)
	return o.outFormat.Print(o.stdout, f)
}

func shootState(o *options) (*cloud_agent.GardenerResponse, error) {
	target := fmt.Sprintf("%s://%s", o.socketNetwork, o.socketAddress)
	o.Logger.Debugf("creating grpc client - target '%s'", target)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
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
