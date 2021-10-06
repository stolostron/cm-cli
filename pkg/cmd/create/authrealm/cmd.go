// Copyright Contributors to the Open Cluster Management project
package authrealm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"

	"github.com/open-cluster-management/cm-cli/pkg/cmd/create/cluster/scenario"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	scenarioDirectory = "create"
)

var valuesTemplatePath = filepath.Join(scenarioDirectory, "values-template.yaml")

var example = `
# Create a cluster
%[1]s create cluster --values values.yaml

# Create a cluster with cluster name overwrite
%[1]s create cluster --cluster mycluster --values values.yaml

# Create a cluster with cluster name overwrite by args
%[1]s create cluster mycluster --values values.yaml
`

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:          "cluster",
		Short:        "Create a cluster",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(c *cobra.Command, args []string) error {
			if !helpers.IsRHACM(cmFlags.KubectlFactory) && !helpers.IsMCE(cmFlags.KubectlFactory) {
				return fmt.Errorf("this command '%s %s' is only available on %s or %s",
					helpers.GetExampleHeader(),
					strings.Join(os.Args[1:], " "),
					helpers.RHACM,
					helpers.MCE)
			}
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.complete(c, args); err != nil {
				return err
			}
			if err := o.validate(); err != nil {
				return err
			}
			if err := o.run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.SetUsageTemplate(clusteradmhelpers.UsageTempate(cmd, scenario.GetScenarioResourcesReader(), valuesTemplatePath))
	cmd.Flags().StringVar(&o.clusterName, "cluster", "", "Name of the cluster")
	cmd.Flags().StringVar(&o.valuesPath, "values", "", "The files containing the values")
	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "The generated resources will be copied in the specified file")
	cmd.Flags().BoolVar(&o.waitAgent, "wait", false, "Wait until the klusterlet agent is installed")
	//Not implemented as it requires to import all addon packages
	// cmd.Flags().BoolVar(&o.waitAddOns, "wait-addons", false, "Wait until the klusterlet agent and the addons are is installed")
	cmd.Flags().IntVar(&o.timeout, "timeout", 180, "Timeout to get the klusterlet agent or addons ready in seconds")

	return cmd
}
