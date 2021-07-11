// Copyright Contributors to the Open Cluster Management project
package get

import (
	"github.com/open-cluster-management/cm-cli/pkg/cmd/get/clusterclaims"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/get/clusterpoolhosts"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/get/clusters"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/get/machinepools"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	clusteradmgettoken "open-cluster-management.io/clusteradm/pkg/cmd/get/token"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get a resource",
	}

	cmd.AddCommand(clusters.NewCmd(cmFlags, streams))
	cmd.AddCommand(machinepools.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusteradmgettoken.NewCmd(clusteradmFlags, streams))
	cmd.AddCommand(clusterpoolhosts.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterclaims.NewCmd(cmFlags, streams))

	return cmd
}
