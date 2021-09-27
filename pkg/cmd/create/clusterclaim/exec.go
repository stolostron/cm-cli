// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"
	"strings"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"k8s.io/apimachinery/pkg/runtime/schema"

	printclusterpoolv1alpha1 "github.com/open-cluster-management/cm-cli/api/cm-cli/v1alpha1"
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
	return nil
}

func (o *Options) run() (err error) {
	cphs, err := clusterpoolhost.GetClusterPoolHosts()
	if err != nil {
		return err
	}

	cph, err := cphs.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	err = cph.CreateClusterClaims(o.ClusterClaims, o.ClusterPool, o.SkipSchedule, o.Timeout, o.CMFlags.DryRun, o.outputFile)
	if err != nil {
		return err
	}

	for _, clusterClaim := range strings.Split(o.ClusterClaims, ",") {
		cc, err := cph.GetClusterClaim(clusterClaim, o.Timeout, o.CMFlags.DryRun, o.GetOptions.PrintFlags)
		if err != nil {
			return err
		}

		cred, err := cph.GetClusterClaimCred(cc, o.WithCredentials)
		if err != nil {
			return err
		}
		cred.GetObjectKind().
			SetGroupVersionKind(
				schema.GroupVersionKind{
					Group:   printclusterpoolv1alpha1.GroupName,
					Kind:    "PrintClusterClaimCredential",
					Version: printclusterpoolv1alpha1.GroupVersion.Version})
		err = helpers.Print(cred, o.GetOptions.PrintFlags)
		if err != nil {
			return err
		}
	}
	return nil
}
