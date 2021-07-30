// Copyright Contributors to the Open Cluster Management project
package clusterpoolhosts

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Get cluster pool hosts
%[1]s get cph

# Get cluster pool hosts
%[1]s get cph -oyaml|json|custom-columns=%[2]s

# Get cluster pool hosts in a raw format
%[1]s get cph --raw
`

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)

	cmd := &cobra.Command{
		Use:          "clusterpoolhosts",
		Aliases:      []string{"clusterpoolhost", "cphs", "cph"},
		Short:        "list the clusterpoolhosts",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader(), clusterpoolhost.ClusterPoolHostsColumns),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return clusterpoolhost.BackupCurrentContexts()
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run())
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return clusterpoolhost.RestoreCurrentContexts()
		},
	}

	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "", "Output format. One of: json|yaml|custom-columns=c1|c2|...")
	return cmd
}
