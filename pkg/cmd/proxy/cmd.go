// Copyright Contributors to the Open Cluster Management project
package proxy

import (
	"fmt"

	"github.com/stolostron/cm-cli/pkg/cmd/proxy/health"
	"github.com/stolostron/cm-cli/pkg/cmd/proxy/kubectl"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "proxy commands",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			rhacmConstraint := ">=2.5.0"
			supported, platform, err := helpers.IsSupportedVersion(cmFlags, false, "", rhacmConstraint, "")
			if err != nil {
				return err
			}
			if platform != helpers.RHACM || !supported {
				version, _, err := helpers.GetVersion(cmFlags, false, "")
				if err != nil {
					return err
				}
				return fmt.Errorf("this command is only valid on RHACM %s current vesion: %s", rhacmConstraint, version)
			}
			return nil
		},
	}

	cmd.AddCommand(health.NewCmd(clusteradmFlags, cmFlags, streams))
	cmd.AddCommand(kubectl.NewCmd(clusteradmFlags, cmFlags, streams))

	return cmd
}
