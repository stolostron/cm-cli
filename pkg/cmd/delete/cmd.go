// Copyright Contributors to the Open Cluster Management project
package delete

import (
	"github.com/stolostron/cm-cli/pkg/cmd/delete/cluster"
	"github.com/stolostron/cm-cli/pkg/cmd/delete/clusterclaim"
	"github.com/stolostron/cm-cli/pkg/cmd/delete/clusterpool"
	"github.com/stolostron/cm-cli/pkg/cmd/delete/clusterpoolhost"
	"github.com/stolostron/cm-cli/pkg/cmd/delete/hypershiftdeployment"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	clusteradmclusterset "open-cluster-management.io/clusteradm/pkg/cmd/delete/clusterset"
	clusteradmclusterwork "open-cluster-management.io/clusteradm/pkg/cmd/delete/work"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a resource",
	}

	cmd.AddCommand(cluster.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterpoolhost.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterclaim.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterpool.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusteradmclusterset.NewCmd(clusteradmFlags, streams))
	cmd.AddCommand(clusteradmclusterwork.NewCmd(clusteradmFlags, streams))
	cmd.AddCommand(hypershiftdeployment.NewCmd(cmFlags, streams))
	return cmd
}
