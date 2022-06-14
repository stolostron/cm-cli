// Copyright Contributors to the Open Cluster Management project
package health

import (
	"fmt"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	clusteradmproxyhealth "open-cluster-management.io/clusteradm/pkg/cmd/proxy/health"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# check health of a managedcluster proxy 
%[1]s proxy health
`

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := clusteradmproxyhealth.NewCmd(clusteradmFlags, streams)
	cmd.Use = "health"
	cmd.Example = fmt.Sprintf(example, clusteradmhelpers.GetExampleHeader())

	return cmd
}
