// Copyright Contributors to the Open Cluster Management project
package unbind

import (
	"github.com/stolostron/cm-cli/pkg/cmd/set/clusterpoolhost"
	"github.com/stolostron/cm-cli/pkg/cmd/unbind/clusterset"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbind",
		Short: "unbind a resource",
	}

	cmd.AddCommand(clusterpoolhost.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterset.NewCmd(clusteradmFlags, cmFlags, streams))

	return cmd
}
