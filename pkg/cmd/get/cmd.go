// Copyright Contributors to the Open Cluster Management project
package get

import (
	"github.com/stolostron/cm-cli/pkg/cmd/get/clusterclaim"
	"github.com/stolostron/cm-cli/pkg/cmd/get/clusterpoolhosts"
	"github.com/stolostron/cm-cli/pkg/cmd/get/clusterpools"
	"github.com/stolostron/cm-cli/pkg/cmd/get/clusters"
	"github.com/stolostron/cm-cli/pkg/cmd/get/config"
	"github.com/stolostron/cm-cli/pkg/cmd/get/credentials"
	"github.com/stolostron/cm-cli/pkg/cmd/get/machinepools"
	"github.com/stolostron/cm-cli/pkg/cmd/get/policies"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	clusteradmaddon "open-cluster-management.io/clusteradm/pkg/cmd/get/addon"
	clusteradmclusterset "open-cluster-management.io/clusteradm/pkg/cmd/get/clusterset"
	clusteradmhubinfo "open-cluster-management.io/clusteradm/pkg/cmd/get/hubinfo"
	clusteradmgettoken "open-cluster-management.io/clusteradm/pkg/cmd/get/token"
	clusteradmwork "open-cluster-management.io/clusteradm/pkg/cmd/get/work"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(f cmdutil.Factory, clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get a resource",
	}

	cmd.AddCommand(clusters.NewCmd(cmFlags, streams))
	cmd.AddCommand(credentials.NewCmd(cmFlags, streams))
	cmd.AddCommand(machinepools.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusteradmgettoken.NewCmd(clusteradmFlags, streams))
	cmd.AddCommand(clusterpoolhosts.NewCmd(f, cmFlags, streams))
	cmd.AddCommand(clusterclaim.NewCmd(f, cmFlags, streams))
	cmd.AddCommand(clusterpools.NewCmd(f, cmFlags, streams))
	cmd.AddCommand(config.NewCmd(clusteradmFlags, cmFlags, streams))
	cmd.AddCommand(policies.NewCmd(f, cmFlags, streams))
	cmd.AddCommand(clusteradmhubinfo.NewCmd(clusteradmFlags, streams))
	cmd.AddCommand(clusteradmaddon.NewCmd(clusteradmFlags, streams))
	cmd.AddCommand(clusteradmclusterset.NewCmd(clusteradmFlags, streams))
	cmd.AddCommand(clusteradmwork.NewCmd(clusteradmFlags, streams))

	return cmd
}
