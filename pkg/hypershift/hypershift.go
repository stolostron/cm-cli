// Copyright Contributors to the Open Cluster Management project
package hypershift

import (
	"context"
	"fmt"

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

	hd := &hypershiftdeploymentv1alpha1.HypershiftDeployment{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(hdus.Items[0].UnstructuredContent(), hd); err != nil {
		return nil, err
	}
	return hd, nil
}
