// Copyright Contributors to the Open Cluster Management project
package clusterset

import (
	"fmt"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	clusteradmclustersetset "open-cluster-management.io/clusteradm/pkg/cmd/clusterset/set"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Set clusters to a clusterset
%[1]s set clusterset clusterset1 --clusters cluster1,cluster2
`

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := clusteradmclustersetset.NewCmd(clusteradmFlags, streams)
	cmd.Use = "clusterset"
	cmd.Example = fmt.Sprintf(example, clusteradmhelpers.GetExampleHeader())
	return cmd
}
