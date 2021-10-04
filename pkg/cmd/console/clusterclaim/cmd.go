// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"

	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	example = `
	# open the console of clusterclaim in current clusterpoolhost
	%[1]s console cc <cluster-name>
	
	# pen the console clusterclaims of a given clusterpoolhost
	%[1]s concole cc <cluster-name> --cph <clusterpoolhosts>
`
)

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {

	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:                   "clusterclaims",
		Aliases:               []string{"cc", "ccs", "clusterclaim"},
		DisableFlagsInUseLine: true,
		Short:                 "Open console of a clusterclaim",
		Example:               fmt.Sprintf(example, helpers.GetExampleHeader()),
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

	o.GetOptions.PrintFlags = get.NewGetPrintFlags()

	o.GetOptions.PrintFlags.AddFlags(cmd)

	cmd.Flags().StringVar(&o.ClusterPoolHost, "cph", "", "The clusterpoolhost to use")
	cmd.Flags().IntVar(&o.Timeout, "timeout", 60, "Timeout to get the cluster claim running")
	cmd.Flags().BoolVar(&o.WithCredentials, "creds", o.WithCredentials, "If set the credentials will be displayed")

	return cmd
}
