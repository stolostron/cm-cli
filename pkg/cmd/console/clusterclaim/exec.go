// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 0 {
		o.ClusterClaim = args[0]
	}
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

	err = cph.OpenClusterClaim(o.ClusterClaim, o.Timeout)
	if err != nil {
		return err
	}

	if o.WithCredentials {
		cc, err := cph.GetClusterClaim(o.ClusterClaim, o.Timeout, o.CMFlags.DryRun, o.GetOptions.PrintFlags)
		if err != nil {
			return err
		}
		return cph.PrintClusterClaimCred(cc, o.GetOptions.PrintFlags, o.WithCredentials)
	}
	return nil

}
