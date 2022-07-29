package config

import (
	"fmt"
	"strings"

	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	addArg = "add"
	delArg = "del"
)

var (
	validArgs = []string{addArg, delArg}
)

func newGardenerCmd(o *gardenerOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gardener [" + strings.Join(validArgs, "|") + "]",
		Short: "Manage gardener projecs in the configuration file.",
		Long:  "Add or delete gardeners project to the configuration file.",
		Args: cobra.MatchAll(
			cobra.MinimumNArgs(1),
			cobra.MaximumNArgs(1),
			cobra.OnlyValidArgs,
		),
		ValidArgs: validArgs,
		PreRunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case addArg:
				return o.validateAdd()
			case delArg:
				return o.validateDel()
			}
			return errors.New(fmt.Sprintf("unsupported argument: %s", args[0]))
		},
		RunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case addArg:
				return runAddGardener(o)
			case delArg:
				return runDelGardener(o)
			}
			return errors.New(fmt.Sprintf("unsupported argument: %s", args[0]))
		},
		Example: `  # Add gardener project
  cloudagent config gardener add --namespace <namespace> --kubeconfig <path>

  # Delete gardener project
  cloudagent config gardener del --namespace <namespace>`,
	}

	cmd.Flags().StringVarP(&o.namespace, "namespace", "n", "", "The name of the gardener project.")
	cmd.Flags().StringVarP(&o.kubeconfig, "kubeconfig", "k", "", "The path of the gardener kubeconfig.")

	return cmd
}

func runAddGardener(o *gardenerOptions) error {
	c, err := o.readConfig(o.configPath)
	if err != nil {
		return errors.Wrap(err, "while reading config file")
	}

	c.GardenerProjects = append(c.GardenerProjects, config.GardenerProject{
		Namespace:      o.namespace,
		KubeconfigPath: o.kubeconfig,
	})

	return o.writeConfig(o.configPath, c)
}

func runDelGardener(o *gardenerOptions) error {
	c, err := o.readConfig(o.configPath)
	if err != nil {
		return errors.Wrap(err, "while reading config from file")
	}

	projects := []config.GardenerProject{}
	for i := range c.GardenerProjects {
		if !isProjectExpected(o.namespace, o.kubeconfig, c.GardenerProjects[i]) {
			projects = append(projects, c.GardenerProjects[i])
		}
	}

	if len(projects) == len(c.GardenerProjects) || len(c.GardenerProjects) == 0 {
		return errors.New("can't find a project which meets all requirements")
	}

	c.GardenerProjects = projects

	err = o.writeConfig(o.configPath, c)
	if err != nil {
		return errors.Wrap(err, "while writing config to file")
	}

	return nil
}

func isProjectExpected(namespace, kubeconfig string, p config.GardenerProject) bool {
	if namespace == "" {
		namespace = p.Namespace
	} else if kubeconfig == "" {
		kubeconfig = p.KubeconfigPath
	}

	if kubeconfig == p.KubeconfigPath && namespace == p.Namespace {
		return true
	}

	return false
}
