// Copyright Contributors to the Open Cluster Management project
package set

import (
	"github.com/stolostron/cm-cli/pkg/cmd/set/clusterpoolhost"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "set a resource",
	}

	cmd.AddCommand(clusterpoolhost.NewCmd(cmFlags, streams))

	return cmd
}
