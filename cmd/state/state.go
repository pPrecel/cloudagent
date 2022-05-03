package state

import (
	"errors"
	"fmt"
	"time"

	"github.com/pPrecel/cloud-agent/internal/output"
	cloud_agent "github.com/pPrecel/cloud-agent/pkg/agent/proto"
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

	cmd.Flags().StringVarP(&o.createdBy, "createdBy", "c", "", "Show clusters created by specific person.")
	cmd.Flags().VarP(output.New(&o.outFormat, "table", "%r/%h/%u/%a", "-/-/-/-"), "output", "o", `Provides format for the output information. 
	
For the 'text' output format you can specifie two more informations by spliting them using '='. The first one would be used as output format and second as error format.

The first one can contains at least on out of four elements where:
- '%r' represents number of running clusters, 
- '%h' represents number of hibernated clusters, 
- '%u' represents number of cluster with unknown status, 
- '%a' represents of all cluster in namespace.

The second one can contains '%e'  which will be replaced with error message.`)
	cmd.Flags().DurationVarP(&o.timeout, "timeout", "t", 2*time.Second, "Provides timeout for the command.")

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
	if s == nil {
		s = &cloud_agent.ShootList{}
	}
	w := o.writer

	o.Logger.Debugf("printing shoots in format '%s'", o.outFormat.String())
	switch o.outFormat.Type() {
	case string(output.JsonType):
		var f []string
		if o.createdBy != "" {
			f = append(f, fmt.Sprintf(`#(annotations.%s=="%s")#`, escapedCreatedByLabel, o.createdBy))
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
		if o.createdBy != "" {
			f = append(f, fmt.Sprintf("#(owner==%s)#", o.createdBy))
		}

		return output.PrintTable(w, tab, f...)
	case string(output.TextType):
		if e == nil && len(s.Shoots) == 0 {
			e = errors.New("empty shoot list")
		}

		if e != nil {
			return output.PrintErrorText(w, output.ErrorOptions{
				Format: o.outFormat.ErrorFormat(),
				Error:  e.Error(),
			})
		} else {
			var f []string
			if o.createdBy != "" {
				f = append(f, fmt.Sprintf(`#(annotations.%s=="%s")#`, escapedCreatedByLabel, o.createdBy))
			}

			return output.PrintText(w, s.Shoots, output.TextOptions{
				Format: o.outFormat.StringFormat(),
				APath:  "#",
				RPath:  `#(condition==1)#|#`,
				HPath:  `#(condition==2)#|#`,
				UPath:  `#(condition==3)#|#`,
			}, f...)
		}
	}

	return nil
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
