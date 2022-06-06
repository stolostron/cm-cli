// Copyright Contributors to the Open Cluster Management project
package hypershift

import (
	"context"
	"fmt"
	"strings"

	"github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	hypershiftdeploymentv1alpha1 "github.com/stolostron/hypershift-deployment-controller/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func GetHypershiftDeployment(clusterName string, cmFlags *genericclioptions.CMFlags) (*hypershiftdeploymentv1alpha1.HypershiftDeployment, error) {
	dynamicClient, err := cmFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return nil, err
	}

	hdus, err := dynamicClient.Resource(helpers.GvrHD).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("hypershift.openshift.io/infra-id=%s", clusterName),
	})
	if err != nil {
		return nil, err
	}

	if len(hdus.Items) == 0 {
		return nil, fmt.Errorf("no hypershiftdeployment found for infra-id: %s", clusterName)
	}

	hd := &hypershiftdeploymentv1alpha1.HypershiftDeployment{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(hdus.Items[0].UnstructuredContent(), hd); err != nil {
		return nil, err
	}
	return hd, nil
}

func DeleteHypershiftDeployments(hypershiftDeployments string, namespace string, cmFlags *genericclioptions.CMFlags) error {
	dynamicClient, err := cmFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}

	for _, hdn := range strings.Split(hypershiftDeployments, ",") {
		hypershiftDeploymentName := strings.TrimSpace(hdn)
		if !cmFlags.DryRun {
			if err := dynamicClient.Resource(helpers.GvrHD).
				Namespace(namespace).Delete(context.TODO(), hypershiftDeploymentName, metav1.DeleteOptions{}); err != nil {
				return err
			}
		} else {
			if _, err := dynamicClient.Resource(helpers.GvrHD).
				Namespace(namespace).Get(context.TODO(), hypershiftDeploymentName, metav1.GetOptions{}); err != nil {
				return err
			}
		}
	}
	return nil
}
