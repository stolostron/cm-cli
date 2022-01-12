// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"

	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Run clusterclaims
%[1]s runcc <clusterclaim_name>[,<clusterclaim_name>...] <options>

# run clusterclaims on a given clusterpoolhost
%[1]s run cc <clusterclaim_name>[,<clusterclaim_name>...] --cph <clusterpoolhost> <options>
`

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:          "clusterclaim",
		Aliases:      []string{"clusterclaims", "cc", "ccs"},
		Short:        "Run clusterclaims",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
			return nil
			// return clusterpoolhost.BackupCurrentContexts()
		},
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
	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "The generated resources will be copied in the specified file")
	cmd.Flags().BoolVar(&o.SkipSchedule, "skip-schedule", false, "Set the hibernation schedule to skip")
	cmd.Flags().IntVar(&o.Timeout, "timeout", 60, "Timeout to get the cluster claim running")
	cmd.Flags().BoolVar(&o.WithCredentials, "creds", o.WithCredentials, "If set the credentials will be displayed")
	return cmd
}
