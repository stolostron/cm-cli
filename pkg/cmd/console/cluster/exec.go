// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"context"
	"fmt"

	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"
	"github.com/stolostron/cm-cli/pkg/helpers"
	"github.com/stolostron/cm-cli/pkg/managedcluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "open-cluster-management.io/api/cluster/v1"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	o.ManagedCluster = args[0]
	return nil
}

func (o *Options) validate() error {
	if o.ManagedCluster == "" {
		return fmt.Errorf("<managed-cluster-name> is missing")
	}
	return nil
}

func (o *Options) run() (err error) {
	dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}

	mcu, err := dynamicClient.Resource(helpers.GvrMC).Get(context.TODO(), o.ManagedCluster, metav1.GetOptions{})
	if err != nil {
		return err
	}
	mc := &clusterv1.ManagedCluster{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(mcu.UnstructuredContent(), mc); err != nil {
		return err
	}

	switch managedcluster.GetClusterType(mc) {
	case managedcluster.HostedType:
		err = o.openHosted(mc)
	case managedcluster.ClusterClaimType:
		err = o.openClusterClaim(mc)
	}
	if err != nil {
		return err
	}
	return nil

}

func (o *Options) openHosted(mc *clusterv1.ManagedCluster) error {
	err := managedcluster.OpenManagedCluster(mc)
	if err != nil {
		return err
	}
	if o.WithCredentials {
		return fmt.Errorf("credentials not available yet")
	}
	return nil
}

func (o *Options) openClusterClaim(mc *clusterv1.ManagedCluster) error {
	cph, err := clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	err = managedcluster.OpenManagedCluster(mc)
	if err != nil {
		return err
	}

	if o.WithCredentials {
		cc, err := cph.GetClusterClaim(o.ManagedCluster, true, o.Timeout, o.CMFlags.DryRun, o.GetOptions.PrintFlags)
		if err != nil {
			return err
		}
		return cph.PrintClusterClaimCred(cc, o.GetOptions.PrintFlags, o.WithCredentials)
	}
	return nil
}
