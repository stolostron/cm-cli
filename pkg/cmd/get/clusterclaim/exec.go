// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"context"
	"fmt"

	printclusterpoolv1alpha1 "github.com/stolostron/cm-cli/api/cm-cli/v1alpha1"
	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"
	"github.com/stolostron/cm-cli/pkg/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"

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
		switch o.KubeConfig {
		case true:
			err = o.getKubeConfig(cph)
		case false:
			err = o.getCC(cph)
		}
	}
	return err

}

func (o *Options) getKubeConfig(cph *clusterpoolhost.ClusterPoolHost) (err error) {
	cc, err := cph.GetClusterClaim(o.ClusterClaim, o.WithCredentials, o.Timeout, o.CMFlags.DryRun, o.GetOptions.PrintFlags)
	if err != nil {
		return err
	}
	cd, err := cph.GetClusterDeployment(cc)
	if err != nil {
		return err
	}
	kubeConfigSecretName := cd.Spec.ClusterMetadata.AdminKubeconfigSecretRef.Name
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return err
	}
	kubeClient, err := kubernetes.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}
	kubeConfigSecret, err := kubeClient.
		CoreV1().
		Secrets(cd.Namespace).
		Get(context.TODO(), kubeConfigSecretName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	fmt.Println(string(kubeConfigSecret.Data["kubeconfig"]))
	return nil
}

func (o *Options) getCC(cph *clusterpoolhost.ClusterPoolHost) (err error) {
	cc, err := cph.GetClusterClaim(o.ClusterClaim, o.WithCredentials, o.Timeout, o.CMFlags.DryRun, o.GetOptions.PrintFlags)
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
		printClusterClaimsList := cph.ConvertToPrintClusterClaimList(clusterClaims, o.Current)
		printClusterClaimLists.Items = append(printClusterClaimLists.Items, printClusterClaimsList.Items...)
	}
	helpers.Print(printClusterClaimLists, o.GetOptions.PrintFlags)
	return nil
}
