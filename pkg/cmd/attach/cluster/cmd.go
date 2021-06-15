// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"fmt"
	"path/filepath"

	"github.com/open-cluster-management/cm-cli/pkg/cmd/attach/cluster/scenario"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Attach a cluster
%[1]s attach cluster --values values.yaml

# Attach a cluster with overwritting the cluster name
%[1]s attach cluster --values values.yaml --cluster mycluster
`

const (
	scenarioDirectory = "attach"
)

var valuesTemplatePath = filepath.Join(scenarioDirectory, "values-template.yaml")
var valuesDefaultPath = filepath.Join(scenarioDirectory, "values-default.yaml")

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)
	// cmd := &cobra.Command{
	// 	Use: "attach",
	// }

	cluster := &cobra.Command{
		Use:          "cluster",
		Short:        "Import a cluster",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(c *cobra.Command, args []string) error {
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
			if !helpers.IsRHACM(cmFlags.KubectlFactory) {
				return fmt.Errorf("this command '%s attach cluster' is only available on RHACM", helpers.GetExampleHeader())
			}
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
	cluster.Flags().StringVar(&o.clusterServer, "cluster-server", "", "cluster server url of the cluster to import")
	cluster.Flags().StringVar(&o.clusterToken, "cluster-token", "", "token to access the cluster to import")
	cluster.Flags().StringVar(&o.clusterKubeConfig, "cluster-kubeconfig", "", "path to the kubeconfig the cluster to import")
	cluster.Flags().StringVar(&o.importFile, "import-file", "", "the file which will contain the import secret for manual import")
	cluster.Flags().StringVar(&o.outputFile, "output-file", "", "The generated resources will be copied in the specified file")

	// cmd.AddCommand(clusters)

	return cluster
}
