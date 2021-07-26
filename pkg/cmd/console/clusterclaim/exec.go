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
	cphs, err := clusterpoolhost.GetClusterPoolHosts()
	if err != nil {
		return err
	}

	currentCph, err := cphs.GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}

	err = clusterpoolhost.OpenClusterClaim(o.ClusterClaim, o.Timeout)
	if err != nil {
		return err
	}

	if len(o.ClusterPoolHost) != 0 {
		if err := cphs.SetActive(currentCph); err != nil {
			return err
		}
	}
	return err

}
