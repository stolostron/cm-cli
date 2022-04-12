// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/helpers"
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
func NewCmd(f cmdutil.Factory, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {

	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:                   "clusterclaims",
		Aliases:               []string{"cc", "ccs", "clusterclaim"},
		DisableFlagsInUseLine: true,
		Short:                 "Display clusterclaims",
		Example:               fmt.Sprintf(example, helpers.GetExampleHeader()),
		// PreRunE: func(cmd *cobra.Command, args []string) error {
		// 	return clusterpoolhost.BackupCurrentContexts()
		// },
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.GetOptions.Complete(f, cmd, []string{"printclusterclaims"}))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run())
		},
		// PostRunE: func(cmd *cobra.Command, args []string) error {
		// 	return clusterpoolhost.RestoreCurrentContexts()
		// },
	}

	o.GetOptions.PrintFlags = get.NewGetPrintFlags()

	o.GetOptions.PrintFlags.AddFlags(cmd)

	cmd.Flags().StringVar(&o.ClusterPoolHost, "cph", "", "The clusterpoolhost to use")
	cmd.Flags().BoolVar(&o.WithCredentials, "creds", o.WithCredentials, "If set the credentials will be displayed")
	cmd.Flags().BoolVarP(&o.AllClusterPoolHosts, "all-cphs", "A", o.AllClusterPoolHosts, "List the clusterclaims across all clusterpoolhosts")
	cmd.Flags().BoolVar(&o.Current, "current", o.Current, "List the clusterclaim which is currently in use")
	cmd.Flags().IntVar(&o.Timeout, "timeout", 60, "Timeout to get the cluster claim running")

	return cmd
}
