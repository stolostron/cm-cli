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
	if len(o.OutputFormat) == 0 {
		o.OutputFormat = helpers.CustomColumnsFormat + clusterpoolhost.ClusterClaimsColumns
	}
	return nil
}

func (o *Options) validate() error {
	if !helpers.IsOutputFormatSupported(o.OutputFormat) {
		return fmt.Errorf("invalid output format %s", helpers.SupportedOutputFormat)
	}
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
	cc, err := clusterpoolhost.GetClusterClaim(o.ClusterClaim, o.Timeout, o.CMFlags.DryRun)
	if err != nil {
		return err
	}
	cred, err := clusterpoolhost.GetClusterClaimCred(cc)
	if err != nil {
		return err
	}
	if o.OutputFormat == helpers.CustomColumnsFormat+clusterpoolhost.ClusterClaimsColumns {
		fmt.Printf("username:    %s\n", cred.User)
		fmt.Printf("password:    %s\n", cred.Password)
		fmt.Printf("basedomain:  %s\n", cred.Basedomain)
		fmt.Printf("api_url:     %s\n", cred.ApiUrl)
		fmt.Printf("console_url: %s\n", cred.ConsoleUrl)
		return nil
	}
	return helpers.Print(cred, o.OutputFormat, o.NoHeaders, nil)

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

	clusterClaimsClaimsP := make([]clusterpoolhost.PrintClusterClaim, 0)
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
		clusterClaimsClaimsP = append(clusterClaimsClaimsP, clusterpoolhost.PrintClusterClaimObj(cphs.ClusterPoolHosts[k], clusterClaims)...)
	}
	helpers.Print(clusterClaimsClaimsP, o.OutputFormat, o.NoHeaders, clusterpoolhost.ConvertClustClaimsForPrint)
	return nil
}
