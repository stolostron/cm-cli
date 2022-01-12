// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"fmt"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Open the console of cluster pool hosts
%[1]s console cph <clusterpoolhost_name>
`

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)

	cmd := &cobra.Command{
		Use:          "clusterpoolhost",
		Aliases:      []string{"clusterpoolhost", "cphs", "cph"},
		Short:        "open the console of a clusterpoolhost",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		// PreRunE: func(cmd *cobra.Command, args []string) error {
		// 	return clusterpoolhost.BackupCurrentContexts()
		// },
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run())
		},
		// PostRunE: func(cmd *cobra.Command, args []string) error {
		// 	return clusterpoolhost.RestoreCurrentContexts()
		// },
	}

	cmd.Flags().StringVar(&o.ClusterPoolHost, "cph", "", "The clusterpoolhost to use")

	return cmd
}
