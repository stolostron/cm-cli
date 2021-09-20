// Copyright Contributors to the Open Cluster Management project
package clusterpoolhosts

import (
	"fmt"

	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Get cluster pool hosts
%[1]s get cph
`

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(f cmdutil.Factory, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)

	cmd := &cobra.Command{
		Use:          "clusterpoolhosts",
		Aliases:      []string{"clusterpoolhost", "cphs", "cph"},
		Short:        "list the clusterpoolhosts",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		// PreRunE: func(cmd *cobra.Command, args []string) error {
		// 	return clusterpoolhost.BackupCurrentContexts()
		// },
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.GetOptions.Complete(f, cmd, []string{"printclusterpoolhost"}))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run())
		},
		// PostRunE: func(cmd *cobra.Command, args []string) error {
		// 	return clusterpoolhost.RestoreCurrentContexts()
		// },
	}

	o.GetOptions.PrintFlags = get.NewGetPrintFlags()

	o.GetOptions.PrintFlags.AddFlags(cmd)
	return cmd
}
