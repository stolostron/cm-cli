// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"context"

	"github.com/stolostron/cm-cli/pkg/genericclioptions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IsIDPInstalled(cmFlags *genericclioptions.CMFlags, skipIDPCheck bool) (isIDPInstalled bool, err error) {
	if skipIDPCheck {
		return true, nil
	}
	f := cmFlags.KubectlFactory
	dynamicClient, err := f.DynamicClient()
	if err != nil {
		panic(err)
	}
	idpConfigs, err := dynamicClient.Resource(GvrIDPConfig).List(context.TODO(), metav1.ListOptions{})
	if err != nil || len(idpConfigs.Items) == 0 {
		return false, err
	}
	return true, nil
}
