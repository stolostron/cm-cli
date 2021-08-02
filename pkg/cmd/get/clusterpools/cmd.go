// Copyright Contributors to the Open Cluster Management project
package clusterpools

import (
	"fmt"

	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Get clusterpool
%[1]s get cp <clusterpool_name>

# Get clusterpool
%[1]s get cp <clusterpool_name> -oyaml|json|custom-columns=%[2]s

# Get clusterpool on a given clusterpoolhost
%[1]s get cp <clusterpool_name> --cph <clusterpoolhost>

# Get clusterpool across all clusterpoolhosts
%[1]s get cp <clusterpool_name> -A
`

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:          "clusterpools",
		Aliases:      []string{"cps", "cp", "clusterpool"},
		Short:        "Get clusterpool",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader(), clusterpoolhost.ClusterPoolsColumns),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
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

	o.PrintFlags = get.NewGetPrintFlags()

	o.PrintFlags.AddFlags(cmd)
	cmd.Flags().StringVar(&o.ClusterPoolHost, "cph", "", "The clusterpoolhost to use")
	cmd.Flags().BoolVarP(&o.AllClusterPoolHosts, "all-cphs", "A", o.AllClusterPoolHosts, "If the requested object does not exist the command will return exit code 0.")

	return cmd
}
