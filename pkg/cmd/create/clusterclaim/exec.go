// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"
	"strings"

	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"
	"github.com/stolostron/cm-cli/pkg/helpers"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("clusterpool name is missing")
	}
	o.ClusterPool = args[0]
	if len(args) < 2 {
		return fmt.Errorf("clusterclaim name is missing")
	}
	o.ClusterClaims = args[1]
	return nil
}

func (o *Options) validate() error {
	if o.Import {
		rhacmConstraint := ">=2.4.0"
		supported, platform, err := helpers.IsSupportedVersion(o.CMFlags, true, o.ClusterPoolHost, rhacmConstraint, "")
		if err != nil {
			return err
		}
		if !supported {
			switch platform {
			case helpers.RHACM:
				return fmt.Errorf("clusterlcaim import is supported only on versions %s", rhacmConstraint)
			case helpers.MCE:
				return fmt.Errorf("clusterlcaim import is supported only on MCE")
			}
		}
	}
	return nil
}

func (o *Options) run() (err error) {
	cph, err := clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	err = cph.CreateClusterClaims(o.ClusterClaims, o.ClusterPool, o.SkipSchedule, o.Import, o.Timeout, o.CMFlags.DryRun, o.outputFile)
	if err != nil {
		return err
	}

	if o.WithCredentials {
		for _, clusterClaim := range strings.Split(o.ClusterClaims, ",") {
			cc, err := cph.GetClusterClaim(clusterClaim, o.Timeout, o.CMFlags.DryRun, o.GetOptions.PrintFlags)
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
