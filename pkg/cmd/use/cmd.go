// Copyright Contributors to the Open Cluster Management project
package use

import (
	"github.com/open-cluster-management/cm-cli/pkg/cmd/use/clusterclaim"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/use/clusterpoolhost"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use",
		Short: "use a resource",
	}

	cmd.AddCommand(clusterclaim.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterpoolhost.NewCmd(cmFlags, streams))

	return cmd
}
