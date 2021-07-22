// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"

	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Use a cluster
%[1]s use cc <cluster_claim_name>

# Use a cluster on a given clusterpoolhosts
%[1]s use cc <cluster_claim_name> --cph <cluster_pool_host_name>
`

// NewCmd provides a cobra command for using a cluster claim
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:          "clusterlclaim",
		Aliases:      []string{"clusterlclaims", "cc", "ccs"},
		Short:        "use clusterclaim change the current context to a cluster claim",
		Long:         "use clusterclaim change the current context to a cluster claim, optionally the cluster pool host can be provided to override the current one",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(c *cobra.Command, args []string) error {
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run())
		},
	}

	cmd.Flags().StringVar(&o.ClusterPoolHost, "cph", "", "The clusterpoolhost to use")
	cmd.Flags().IntVar(&o.Timeout, "timeout", 60, "Timeout to wait the cluster claim running")

	return cmd
}
