// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"

	printclusterpoolv1alpha1 "github.com/open-cluster-management/cm-cli/api/cm-cli/v1alpha1"
	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"k8s.io/apimachinery/pkg/runtime/schema"

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

	cph, err := clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	if len(o.ClusterClaim) == 0 {
		err = o.getCCS(cphs)
	} else {
		err = o.getCC(cph)
	}
	return err

}

func (o *Options) getCC(cph *clusterpoolhost.ClusterPoolHost) (err error) {
	cc, err := cph.GetClusterClaim(o.ClusterClaim, o.Timeout, o.CMFlags.DryRun, o.GetOptions.PrintFlags)
	if err != nil {
		return err
	}
	return cph.PrintClusterClaimCred(cc, o.GetOptions.PrintFlags, o.WithCredentials)
}

func (o *Options) getCCS(cphs *clusterpoolhost.ClusterPoolHosts) (err error) {

	if !o.AllClusterPoolHosts {
		cph, err := clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
		if err != nil {
			return err
		}

		cphs = &clusterpoolhost.ClusterPoolHosts{
			ClusterPoolHosts: map[string]*clusterpoolhost.ClusterPoolHost{
				cph.Name: cph,
			},
		}
	}

	printClusterClaimLists := &printclusterpoolv1alpha1.PrintClusterClaimList{}
	printClusterClaimLists.GetObjectKind().
		SetGroupVersionKind(
			schema.GroupVersionKind{
				Group:   printclusterpoolv1alpha1.GroupName,
				Kind:    "PrintClusterClaim",
				Version: printclusterpoolv1alpha1.GroupVersion.Version})
	for _, cph := range cphs.ClusterPoolHosts {
		clusterClaims, err := cph.GetClusterClaims(o.CMFlags.DryRun)
		if err != nil {
			fmt.Printf("Error while retrieving clusterclaims from %s\n", cph.Name)
			continue
		}
		printClusterClaimsList := cph.ConvertToPrintClusterClaimList(clusterClaims)
		printClusterClaimLists.Items = append(printClusterClaimLists.Items, printClusterClaimsList.Items...)
	}
	helpers.Print(printClusterClaimLists, o.GetOptions.PrintFlags)
	return nil
}
