// Copyright Contributors to the Open Cluster Management project
package hosted

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
	# open the console of a hosted cluster in current clusterpoolhost
	%[1]s console hosted <hosted-cluster-name> <hosting-cluster-name>
	
	# open the console of a hosted cluster of a given clusterpoolhost
	%[1]s concole hosted <hosted-cluster-name> <hosting-cluster-name> --cph <clusterpoolhosts>
`
)

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {

	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:                   "hosteds",
		Aliases:               []string{"hd", "hds", "hosted"},
		DisableFlagsInUseLine: true,
		Short:                 "Open console of a hosted cluster",
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

	cmd.Flags().IntVar(&o.Timeout, "timeout", 60, "Timeout to get the cluster claim running")
	cmd.Flags().BoolVar(&o.WithCredentials, "creds", o.WithCredentials, "If set the credentials will be displayed")

	return cmd
}
