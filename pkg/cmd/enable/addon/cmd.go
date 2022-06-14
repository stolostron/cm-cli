// Copyright Contributors to the Open Cluster Management project
package addon

import (
	"fmt"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	clusteradmaddon "open-cluster-management.io/clusteradm/pkg/cmd/addon/enable"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# enable addons on the hub
%[1]s enable addon 
`

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := clusteradmaddon.NewCmd(clusteradmFlags, streams)
	cmd.Use = "addon"
	cmd.Example = fmt.Sprintf(example, clusteradmhelpers.GetExampleHeader())

	return cmd
}
