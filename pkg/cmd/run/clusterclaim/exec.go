// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"
	"strings"

	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"

	"github.com/spf13/cobra"
)

var scheduleSkip string

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("clusterclaim name is missing")
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
	if cmd.Flags().Lookup("skip-schedule").Changed {
		fmt.Printf("Warninfg: skip-schedule is deprecated, please use hibernate-schedule-on")
		scheduleSkip = "skip"
	}
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

	err = cph.RunClusterClaims(o.ClusterClaims, scheduleSkip, o.Timeout, o.CMFlags.DryRun, o.outputFile, o.GetOptions.PrintFlags)
	if err != nil {
		return err
	}

	if o.WithCredentials {
		for _, clusterClaim := range strings.Split(o.ClusterClaims, ",") {
			cc, err := cph.GetClusterClaim(clusterClaim, false, o.Timeout, o.CMFlags.DryRun, o.GetOptions.PrintFlags)
			if err != nil {
				return err
			}
			err = cph.PrintClusterClaimCred(cc, o.GetOptions.PrintFlags, o.WithCredentials)
			if err != nil {
				return err
			}
		}
	}
	return nil

}
