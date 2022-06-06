package config

import (
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	addArg = "add"
	delArg = "del"
)

func newGardenerCmd(o *gardenerOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gardener",
		Short: "Manage gardener projecs in the configuration file.",
		Long:  "Add or delete gardeners project to the configuration file.",
		Args: cobra.MatchAll(
			cobra.MinimumNArgs(1),
			cobra.MaximumNArgs(1),
		),
		ValidArgs: []string{"add", "del"},
		PreRunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case addArg:
				return o.validateAdd()
			case delArg:
				return o.validateDel()
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			switch args[0] {
			case addArg:
				return runAddGardener(o)
			case delArg:
				return runDelGardener(o)
			}
			return nil
		},
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

	if len(projects) == len(c.GardenerProjects) {
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
