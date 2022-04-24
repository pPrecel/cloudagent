package generate

import (
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return o.validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(o)
		},
	}

	cmd.Flags().StringVarP(&o.kubeconfigPath, "kubeconfigPath", "k", "", "Provides path to kubeconfig.")
	cmd.Flags().StringVarP(&o.namespace, "namespace", "n", "", "Provides gardener namespace.")
	cmd.Flags().StringVarP(&o.cronSpec, "cronSpec", "c", "@every 60s", "Provices spec for cron configuration.")
	cmd.PersistentFlags().BoolVar(&o.agentVerbose, "agentVerbose", false, "Displays details of actions triggered by the command.")

	return cmd
}

func run(o *options) error {
	o.Logger.Debug("starting command")

	args := []string{}
	args = append(args, "--kubeconfigPath", o.kubeconfigPath)
	args = append(args, "--namespace", o.namespace)

	if o.cronSpec != "" {
		args = append(args, "--cronSpec", o.cronSpec)
	}

	if o.agentVerbose {
		args = append(args, "--verbose")
	}

	procPath, err := o.executable()
	if err != nil {
		return err
	}
	o.Logger.Debugf("main process found in path: \"%s\"", procPath)

	body := darwin.PlistBody(procPath, args)

	o.stdout.Write(body)
	return nil
}
