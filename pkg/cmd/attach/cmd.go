// Copyright Contributors to the Open Cluster Management project
package attach

import (
	"github.com/stolostron/cm-cli/pkg/cmd/attach/cluster"
	"github.com/stolostron/cm-cli/pkg/cmd/attach/clusterclaim"
	"github.com/stolostron/cm-cli/pkg/cmd/attach/hostedcluster"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach",
		Short: "attach a resource",
	}

	cmd.AddCommand(cluster.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterclaim.NewCmd(cmFlags, streams))
	cmd.AddCommand(hostedcluster.NewCmd(cmFlags, streams))

	return cmd
}
