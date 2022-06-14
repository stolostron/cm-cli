// Copyright Contributors to the Open Cluster Management project
package kubectl

import (
	"fmt"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	clusteradmproxykubectl "open-cluster-management.io/clusteradm/pkg/cmd/proxy/kubectl"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# run a kubectl cmd on a managedcluster
%[1]s proxy kubectl cluster-name
`

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := clusteradmproxykubectl.NewCmd(clusteradmFlags, streams)
	cmd.Use = "kubectl"
	cmd.Example = fmt.Sprintf(example, clusteradmhelpers.GetExampleHeader())

	return cmd
}
