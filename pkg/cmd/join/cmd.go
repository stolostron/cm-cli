// Copyright Contributors to the Open Cluster Management project
package join

import (
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	clusteradmjoin "open-cluster-management.io/clusteradm/pkg/cmd/join"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := clusteradmjoin.NewCmd(clusteradmFlags, streams)

	return cmd
}
