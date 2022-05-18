// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"fmt"

	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# set clusters
%[1]s set cluster <cluster_name>[,<clusterc_name>...] <options>
`

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:          "cluster",
		Aliases:      []string{"clusters"},
		Short:        "set clusters",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			clusteradmhelpers.DryRunMessage(cmFlags.DryRun)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate(cmd))
			cmdutil.CheckErr(o.run())
		},
	}
	cmd.Flags().BoolVar(&o.HibernateScheduleOn, "hibernate-schedule-on", false, "Set the hibernation schedule to on")
	cmd.Flags().BoolVar(&o.HibernateScheduleOff, "hibernate-schedule-off", false, "Set the hibernation schedule to off")
	cmd.Flags().BoolVar(&o.HibernateScheduleForce, "hibernate-schedule-force", false, "Force the hibernate setting")

	return cmd
}
