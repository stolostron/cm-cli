// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-cluster-management/cm-cli/pkg/cmd/scale/cluster/scenario"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Scale a cluster
%[1]s scale cluster --cluster clustername --machinepool poolname --replicas 4
`

const (
	scenarioDirectory = "scale"
)

var valuesTemplatePath = filepath.Join(scenarioDirectory, "values-template.yaml")
var valuesDefaultPath = filepath.Join(scenarioDirectory, "values-default.yaml")

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)

	cluster := &cobra.Command{
		Use:          "cluster",
		Short:        "Scale a cluster",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(c *cobra.Command, args []string) error {
			if !helpers.IsRHACM(cmFlags.KubectlFactory) {
				return fmt.Errorf("this command '%s %s' is only available on RHACM", helpers.GetExampleHeader(), strings.Join(os.Args[1:], " "))
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

	cluster.SetUsageTemplate(clusteradmhelpers.UsageTempate(cluster, scenario.GetScenarioResourcesReader(), valuesTemplatePath))
	cluster.Flags().StringVar(&o.valuesPath, "values", "", "The files containing the values")
	cluster.Flags().StringVar(&o.clusterName, "cluster", "", "Name of the cluster")
	cluster.Flags().IntVar(&o.replicas, "replicas", 3, "number of workers for the pool")

	return cluster
}
