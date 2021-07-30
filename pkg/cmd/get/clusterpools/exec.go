// Copyright Contributors to the Open Cluster Management project
package clusterpools

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(o.OutputFormat) == 0 {
		o.OutputFormat = helpers.CustomColumnsFormat + clusterpoolhost.ClusterPoolsColumns
	}
	return nil
}

func (o *Options) validate() error {
	if o.ClusterPoolHost != "" && o.AllClusterPoolHosts {
		return fmt.Errorf("clusterpoolhost and all-cphs are imcompatible")
	}
	if !helpers.IsOutputFormatSupported(o.OutputFormat) {
		return fmt.Errorf("invalid output format %s", helpers.SupportedOutputFormat)
	}
	return nil
}

func (o *Options) run() (err error) {
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

	clusterPoolHostsCPs := make([]clusterpoolhost.PrintClusterPool, 0)
	for k := range cphs.ClusterPoolHosts {
		err = allcphs.SetActive(allcphs.ClusterPoolHosts[k])
		if err != nil {
			return err
		}
		clusterPools, err := clusterpoolhost.GetClusterPools(o.AllClusterPoolHosts, o.CMFlags.DryRun)
		if err != nil {
			fmt.Printf("Error while retrieving clusterpools from %s\n", cphs.ClusterPoolHosts[k].Name)
			continue
		}
		clusterPoolHostsCPs = append(clusterPoolHostsCPs, clusterpoolhost.PrintClusterPoolObj(cphs.ClusterPoolHosts[k], clusterPools)...)
	}
	helpers.Print(clusterPoolHostsCPs, o.OutputFormat, o.NoHeaders, clusterpoolhost.ConvertClusterPoolsForPrint)
	return allcphs.SetActive(currentCph)
}
