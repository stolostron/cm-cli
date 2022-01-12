// Copyright Contributors to the Open Cluster Management project
package with

import (
	"github.com/stolostron/cm-cli/pkg/cmd/with/clusterclaim"
	// "github.com/stolostron/cm-cli/pkg/cmd/with/clusterpoolhost"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmd
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "with",
		Short: "execute a command on a specific cluster",
	}

	cmd.AddCommand(clusterclaim.NewCmd(cmFlags, streams))
	// cmd.AddCommand(clusterpoolhost.NewCmd(cmFlags, streams))

	return cmd
}
