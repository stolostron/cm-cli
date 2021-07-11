// Copyright Contributors to the Open Cluster Management project
package clusterclaims

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 0 {
		o.ClusterPoolHost = args[0]
	}
	return nil
}

func (o *Options) validate() error {
	if o.ClusterPoolHost != "" && o.AllClusterPoolHosts {
		return fmt.Errorf("clusterpoolhost and all-cphs are imcompatible")
	}
	return nil
}

func (o *Options) run() (err error) {
	if err != nil {
		return err
	}
	var cphs, allcphs *clusterpoolhost.ClusterPoolHosts
	allcphs, err = clusterpoolhost.GetClusterPoolHosts()
	if err != nil {
		return err
	}
	currentCph, err := allcphs.GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}

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

	for k := range cphs.ClusterPoolHosts {
		err = allcphs.SetActive(allcphs.ClusterPoolHosts[k])
		if err != nil {
			return err
		}
		err = clusterpoolhost.GetClusterClaims(o.AllClusterPoolHosts, o.CMFlags.DryRun, o.outputFile)
		if err != nil {
			return err
		}
	}
	if len(o.ClusterPoolHost) != 0 {
		return allcphs.SetActive(currentCph)
	}
	return nil
}
