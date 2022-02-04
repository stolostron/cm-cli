// Copyright Contributors to the Open Cluster Management project
package install

import (
	"github.com/stolostron/cm-cli/pkg/cmd/install/acm"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command to install acm
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "install a product",
	}

	cmd.AddCommand(acm.NewCmd(cmFlags, streams))

	return cmd
}
