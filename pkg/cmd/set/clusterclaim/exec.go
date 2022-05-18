// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"

	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"

	"github.com/spf13/cobra"
)

var scheduleSkip string

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("cluster names are missing")
	}
	o.ClusterClaims = args[0]
	if cmd.Flags().Lookup("hibernate-schedule-on").Changed {
		scheduleSkip = "true"
	}
	if cmd.Flags().Lookup("hibernate-schedule-off").Changed {
		scheduleSkip = "skip"
	}

	return nil
}

func (o *Options) validate(cmd *cobra.Command) error {
	if cmd.Flags().Lookup("hibernate-schedule-on").Changed &&
		cmd.Flags().Lookup("hibernate-schedule-off").Changed {
		return fmt.Errorf("flags hibernate-schedule-on and hibernate-schedule-off are mutually exclusif")
	}
	return nil
}

func (o *Options) run() (err error) {

	cph, err := clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	return cph.SetHibernateScheduleClusterClaims(o.ClusterClaims, scheduleSkip, o.CMFlags.DryRun)
}
