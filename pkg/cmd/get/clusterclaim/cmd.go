// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	example = `
	# get a clusterclaim in current clusterpoolhost
	%[1]s get cc <clusterclaim_name> 
	
	# get clusterclaims on a specific clusterpoolhost
	%[1]s get cc  <clusterclaim_name> --cph <clusterpoolhosts>
	
	# get clusterclaims across all clusterpoolhosts
	%[1]s get cc -A`
)

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {

	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:                   "clusterclaims",
		Aliases:               []string{"cc", "ccs", "clusterclaim"},
		DisableFlagsInUseLine: true,
		Short:                 "Display clusterclaims",
		Example:               fmt.Sprintf(example, helpers.GetExampleHeader(), clusterpoolhost.ClusterClaimsColumns),
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

	o.PrintFlags = get.NewGetPrintFlags()

	o.PrintFlags.AddFlags(cmd)

	cmd.Flags().StringVar(&o.ClusterPoolHost, "cph", "", "The clusterpoolhost to use")
	cmd.Flags().BoolVarP(&o.AllClusterPoolHosts, "all-cphs", "A", o.AllClusterPoolHosts, "If the requested object does not exist the command will return exit code 0.")
	cmd.Flags().IntVar(&o.Timeout, "timeout", 60, "Timeout to get the cluster claim running")

	return cmd
}
