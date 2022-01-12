// Copyright Contributors to the Open Cluster Management project
package scale

import (
	"github.com/stolostron/cm-cli/pkg/cmd/scale/cluster"
	"github.com/stolostron/cm-cli/pkg/cmd/scale/clusterpool"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scale",
		Short: "scale a resource",
	}

	cmd.AddCommand(cluster.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterpool.NewCmd(cmFlags, streams))

	return cmd
}
