// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"fmt"

	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Initialize a clusterpool management cluster
%[1]s create clusterpoolhost <clusterpool_name>
`

// var valuesDefaultPath = filepath.Join(scenarioDirectory, "values-default.yaml")

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)

	cmd := &cobra.Command{
		Use:          "clusterpoolhost",
		Aliases:      []string{"cph"},
		Short:        "Initialize a clusterpool management cluster",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(c *cobra.Command, args []string) error {
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

	cmd.Flags().StringVar(&o.ClusterPoolHost.APIServer, "api-server", "", "The API address of the cluster where your 'ClusterPools' are defined. Also referred to as the 'ClusterPool host'")
	cmd.Flags().StringVar(&o.ClusterPoolHost.Console, "console", "", "The URL of the OpenShift console for the ClusterPool host")
	cmd.Flags().StringVar(&o.ClusterPoolHost.Group, "group", "", "Name of a 'Group' ('user.openshift.io/v1') that should be added to each 'ClusterClaim' for team access")
	cmd.Flags().StringVarP(&o.ClusterPoolHost.Namespace, "namespace", "n", "", "Namespace where 'ClusterPools' are defined")
	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "The generated resources will be copied in the specified file")
	cmd.Flags().BoolVar(&o.force, "force", false, "If set and the cluster pool host already exists, it will be overwritten")
	return cmd
}
