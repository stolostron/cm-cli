// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"
	"strings"

	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("clusterclaim name is missing")
	}
	o.ClusterClaims = args[0]
	return nil
}

func (o *Options) validate() error {
	return nil
}

func (o *Options) run() (err error) {
	cph, err := clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	err = cph.RunClusterClaims(o.ClusterClaims, o.SkipSchedule, o.Timeout, o.CMFlags.DryRun, o.outputFile, o.GetOptions.PrintFlags)
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
