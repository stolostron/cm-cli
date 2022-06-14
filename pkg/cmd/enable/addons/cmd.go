// Copyright Contributors to the Open Cluster Management project
package addons

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stolostron/cm-cli/pkg/cmd/enable/addons/scenario"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Enable addons on a cluster
%[1]s enable addons --values values.yaml

# Attach a cluster with overwritting the cluster name
%[1]s enable addons --values values.yaml --cluster mycluster
`

const (
	scenarioDirectory = "addons"
)

var valuesTemplatePath = filepath.Join(scenarioDirectory, "values-template.yaml")
var valuesDefaultPath = filepath.Join(scenarioDirectory, "values-default.yaml")

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)

	cluster := &cobra.Command{
		Use:          "addons",
		Short:        "enable addons on managed cluster",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(c *cobra.Command, args []string) error {
			if !helpers.IsRHACM(o.CMFlags) {
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
	cluster.Flags().StringVar(&o.outputFile, "output-file", "", "The generated resources will be copied in the specified file")

	return cluster
}
