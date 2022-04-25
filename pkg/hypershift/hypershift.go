// Copyright Contributors to the Open Cluster Management project
package hypershift

import (
	"context"
	"fmt"

	"github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubectl/pkg/cmd/get"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
)

const (
	ConsoleURLClusterClaim string = "consoleurl.cluster.open-cluster-management.io"
)

func OpenHosted(cmFlags *genericclioptions.CMFlags, clusterName string, timeout int, printFlags *get.PrintFlags) error {
	dynamicClient, err := cmFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}
	mcu, err := dynamicClient.Resource(helpers.GvrMC).Get(context.TODO(), clusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	mc := &clusterv1.ManagedCluster{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(mcu.UnstructuredContent(), mc); err != nil {
		return err
	}

	var consoleURL string
	for _, c := range mc.Status.ClusterClaims {
		if c.Name == ConsoleURLClusterClaim {
			consoleURL = c.Value
		}
	}
	if consoleURL == "" {
		return fmt.Errorf("console url not found for cluster %s", clusterName)
	}

	return helpers.Openbrowser(consoleURL)
}
