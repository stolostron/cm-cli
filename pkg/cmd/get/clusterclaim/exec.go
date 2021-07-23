// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"

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
	if len(o.ClusterClaim) == 0 {
		err = o.getCCS(cphs)
	} else {
		err = o.getCC(cphs)
	}
	if len(o.ClusterPoolHost) != 0 {
		if err := cphs.SetActive(currentCph); err != nil {
			return err
		}
	}
	return err

}

func (o *Options) getCC(cphs *clusterpoolhost.ClusterPoolHosts) (err error) {
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
	return clusterpoolhost.GetClusterClaim(o.ClusterClaim, o.Timeout, o.CMFlags.DryRun)
}

func (o *Options) getCCS(allcphs *clusterpoolhost.ClusterPoolHosts) (err error) {
	var cphs *clusterpoolhost.ClusterPoolHosts

	if o.AllClusterPoolHosts {
		cphs, err = clusterpoolhost.GetClusterPoolHosts()
		if err != nil {
			return err
		}
	} else {
		var cph *clusterpoolhost.ClusterPoolHost
		if o.ClusterPoolHost != "" {
			cph, err = clusterpoolhost.GetClusterPoolHost(o.ClusterPoolHost)
		} else {
			cph, err = clusterpoolhost.GetCurrentClusterPoolHost()
		}
		if err != nil {
			return err
		}
		cphs = &clusterpoolhost.ClusterPoolHosts{
			ClusterPoolHosts: map[string]*clusterpoolhost.ClusterPoolHost{
				cph.Name: cph,
			},
		}
	}

	allLines := make([]string, 0)
	for k := range cphs.ClusterPoolHosts {
		err = allcphs.SetActive(allcphs.ClusterPoolHosts[k])
		if err != nil {
			return err
		}
		clusterClaims, err := clusterpoolhost.GetClusterClaims(o.CMFlags.DryRun)
		if err != nil {
			fmt.Printf("Error while retrieving clusterclaims from %s\n", cphs.ClusterPoolHosts[k].Name)
			continue
		}
		allLines = append(allLines, clusterpoolhost.SprintClusterClaims(cphs.ClusterPoolHosts[k], "\t", clusterClaims)...)

	}
	helpers.PrintLines(allLines, "\t")
	return nil
}
