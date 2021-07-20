// Copyright Contributors to the Open Cluster Management project
package clusterpool

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("clusterpool name is missing")
	}
	o.ClusterPools = args[0]
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

	err = o.deleteClusterClaims(cphs)

	if len(o.ClusterPoolHost) != 0 {
		if err := cphs.SetActive(currentCph); err != nil {
			return err
		}
	}
	return err
}

func (o *Options) deleteClusterClaims(cphs *clusterpoolhost.ClusterPoolHosts) (err error) {
	if len(o.ClusterPoolHost) != 0 {
		cph, err := cphs.GetClusterPoolHost(o.ClusterPoolHost)
		if err != nil {
			return err
		}

		err = cphs.SetActive(cph)
		if err != nil {
			return err
		}
	}
	return clusterpoolhost.DeleteClusterPools(o.ClusterPools, o.CMFlags.DryRun, o.outputFile)

}
