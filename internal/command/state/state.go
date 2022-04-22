package state

import (
	"fmt"
	"os"
	"time"

	"github.com/pPrecel/cloud-agent/internal/agent"
	cloud_agent "github.com/pPrecel/cloud-agent/internal/agent/proto"
	"github.com/pPrecel/cloud-agent/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	createdByLabel        = `gardener.cloud/created-by`
	escapedCreatedByLabel = `gardener\.cloud/created-by`
)

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "state",
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(o)
		},
	}

	cmd.Flags().StringVarP(&o.CreatedBy, "createdBy", "c", "", "Show clusters created by specific person.")
	cmd.Flags().VarP(output.New(&o.OutFormat, "table", "Shoots: %r/%h/%u/%a", "Error: %e"), "output", "o", `Provides format for the output information. 
	
For the 'text' output format you can specifie two more informations by spliting them using '='. The first one would be used as output format and second as error format.

The first one can contains at least on out of four elements where:
- '%r' represents number of running clusters, 
- '%h' represents number of hibernated clusters, 
- '%u' represents number of cluster with unknown status, 
- '%a' represents of all cluster in namespace.

The second one can contains '%e'  which will be replaced with error message.`)
	cmd.Flags().DurationVarP(&o.Timeout, "timeout", "t", 2*time.Second, "Provides timeout for the command.")

	return cmd
}

func run(o *options) error {
	o.Logger.Debugf("getting shoots")
	list, err := shootState(o)

	o.Logger.Debugf("received: %+v, error: %v", list, err)

	return printOutput(o, list, err)
}

type tableFormat struct {
	Name      string `header:"Name" json:"name"`
	Owner     string `header:"Created By" json:"owner"`
	Condition string `header:"Condition" json:"condition"`
}

func printOutput(o *options, s *cloud_agent.ShootList, e error) error {
	w := os.Stdout

	o.Logger.Debugf("printing shoots in format '%s'", o.OutFormat.String())
	switch o.OutFormat.Type() {
	case string(output.JsonType):
		var f []string
		if o.CreatedBy != "" {
			f = append(f, fmt.Sprintf(`#(annotations.%s=="%s")#`, escapedCreatedByLabel, o.CreatedBy))
		}

		return output.PrintJson(w, s.Shoots, f...)
	case string(output.TableType):
		tab := []tableFormat{}

		for i := range s.Shoots {
			tab = append(tab, tableFormat{
				Name:      s.Shoots[i].Name,
				Owner:     s.Shoots[i].Annotations[createdByLabel],
				Condition: s.Shoots[i].Condition.String(),
			})
		}

		var f []string
		if o.CreatedBy != "" {
			f = append(f, fmt.Sprintf("#(owner==%s)#", o.CreatedBy))
		}

		return output.PrintTable(w, tab, f...)
	case string(output.TextType):
		if e != nil {
			return output.PrintErrorText(w, output.ErrorOptions{
				Format: o.OutFormat.ErrorFormat(),
				Error:  e.Error(),
			})
		} else {
			var f []string
			if o.CreatedBy != "" {
				f = append(f, fmt.Sprintf(`#(annotations.%s=="%s")#`, escapedCreatedByLabel, o.CreatedBy))
			}

			return output.PrintText(w, s.Shoots, output.TextOptions{
				Format: o.OutFormat.StringFormat(),
				APath:  "#",
				RPath:  `#(condition==0)#|#`,
				HPath:  `#(condition==1)#|#`,
				UPath:  `#(condition==2)#|#`,
			}, f...)
		}
	}

	return nil
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
