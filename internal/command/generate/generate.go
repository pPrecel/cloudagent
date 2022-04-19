package generate

import (
	"fmt"
	"os"

	"github.com/pPrecel/cloud-agent/internal/darwin"
	"github.com/spf13/cobra"
)

func NewCmd(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "generate [plist]",
		Short:                 "Generate some system-oriented utils.",
		Long:                  "Use this command to generate a system-oriented utils it would help you in communication between the cloudagent and other tools.",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"plist"},
		Args:                  cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(o)
		},
	}

	cmd.Flags().StringVarP(&o.KubeconfigPath, "kubeconfigPath", "k", "", "Provides path to kubeconfig.")
	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", "", "Provides gardener namespace.")
	cmd.Flags().StringVarP(&o.CronSpec, "cronSpec", "c", "@every 60s", "Provices spec for cron configuration.")
	cmd.PersistentFlags().BoolVar(&o.AgentVerbose, "agentVerbose", false, "Displays details of actions triggered by the command.")

	return cmd
}

func run(o *options) error {
	o.Logger.Debug("starting command")
	args := []string{}
	if o.KubeconfigPath != "" {
		args = append(args, "--kubeconfigPath", o.KubeconfigPath)
	}
	if o.Namespace != "" {
		args = append(args, "--namespace", o.Namespace)
	}
	if o.CronSpec != "" {
		args = append(args, "--cronSpec", o.CronSpec)
	}
	if o.AgentVerbose {
		args = append(args, "--verbose")
	}

	procPath, err := os.Executable()
	if err != nil {
		return err
	}
	o.Logger.Debugf("main process found in path: \"%s\"", procPath)

	body := darwin.PlistBody(procPath, args)

	fmt.Printf("%s", body)
	return nil
}
