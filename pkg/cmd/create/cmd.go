// Copyright Contributors to the Open Cluster Management project
package create

import (
	"github.com/stolostron/cm-cli/pkg/cmd/create/authrealm"
	"github.com/stolostron/cm-cli/pkg/cmd/create/cluster"
	"github.com/stolostron/cm-cli/pkg/cmd/create/clusterclaim"
	"github.com/stolostron/cm-cli/pkg/cmd/create/clusterpool"
	"github.com/stolostron/cm-cli/pkg/cmd/create/clusterpoolhost"
	"github.com/stolostron/cm-cli/pkg/cmd/create/hypershiftdeployment"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	clusteradmclusterset "open-cluster-management.io/clusteradm/pkg/cmd/create/clusterset"
	clusteradmwork "open-cluster-management.io/clusteradm/pkg/cmd/create/work"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a resource",
	}

	cmd.AddCommand(authrealm.NewCmd(cmFlags, streams))
	cmd.AddCommand(cluster.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterpoolhost.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterclaim.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterpool.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusteradmclusterset.NewCmd(clusteradmFlags, streams))
	cmd.AddCommand(clusteradmwork.NewCmd(clusteradmFlags, streams))
	cmd.AddCommand(hypershiftdeployment.NewCmd(cmFlags, streams))

	return cmd
}
