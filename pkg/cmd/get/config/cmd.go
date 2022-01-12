// Copyright Contributors to the Open Cluster Management project
package config

import (
	"github.com/stolostron/cm-cli/pkg/cmd/get/config/cluster"
	"github.com/stolostron/cm-cli/pkg/cmd/get/config/clusterpool"
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "get the config of a resource",
	}

	cmd.AddCommand(cluster.NewCmd(cmFlags, streams))
	cmd.AddCommand(clusterpool.NewCmd(cmFlags, streams))

	return cmd
}
