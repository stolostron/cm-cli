// Copyright Contributors to the Open Cluster Management project
package detach

import (
	"github.com/open-cluster-management/cm-cli/pkg/cmd/detach/cluster"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach",
		Short: "detach a resources",
	}

	cmd.AddCommand(cluster.NewCmd(cmFlags, streams))

	return cmd
}
