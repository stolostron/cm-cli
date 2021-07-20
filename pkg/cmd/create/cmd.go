// Copyright Contributors to the Open Cluster Management project
package create

import (
	"github.com/open-cluster-management/cm-cli/pkg/cmd/create/cluster"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/create/clusterclaim"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/create/clusterpool"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/create/clusterpoolhost"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a resource",
	}

	cmd.AddCommand(cluster.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterpoolhost.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterclaim.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterpool.NewCmd(cmFlags, streams))

	return cmd
}
