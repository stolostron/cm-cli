// Copyright Contributors to the Open Cluster Management project
package clusterset

import (
	"fmt"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	clusteradmunbind "open-cluster-management.io/clusteradm/pkg/cmd/clusterset/unbind"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# UnBind a clusterset from a namespace
%[1]s unbind clusterset clusterset1 --namespace default
`

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := clusteradmunbind.NewCmd(clusteradmFlags, streams)
	cmd.Use = "clusterset"
	cmd.Example = fmt.Sprintf(example, clusteradmhelpers.GetExampleHeader())

	return cmd
}
