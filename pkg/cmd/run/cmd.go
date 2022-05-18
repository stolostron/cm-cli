// Copyright Contributors to the Open Cluster Management project
package run

import (
	"github.com/stolostron/cm-cli/pkg/cmd/run/cluster"
	"github.com/stolostron/cm-cli/pkg/cmd/run/clusterclaim"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run a resource",
	}

	cmd.AddCommand(clusterclaim.NewCmd(cmFlags, streams))
	cmd.AddCommand(cluster.NewCmd(cmFlags, streams))

	return cmd
}
