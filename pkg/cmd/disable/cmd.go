// Copyright Contributors to the Open Cluster Management project
package disable

import (
	"github.com/stolostron/cm-cli/pkg/cmd/disable/components"
	"github.com/stolostron/cm-cli/pkg/cmd/enable/addon"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "disable a feature",
	}

	cmd.AddCommand(components.NewCmd(cmFlags, streams))
	cmd.AddCommand(addon.NewCmd(clusteradmFlags, cmFlags, streams))

	return cmd
}
