// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"context"
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IsSupported(cmFlags *genericclioptions.CMFlags) (isSupported bool, err error) {
	return cmFlags.SkipServerCheck || IsRHACM(cmFlags) || IsMCE(cmFlags), err
}

func IsRHACM(cmFlags *genericclioptions.CMFlags) bool {
	cms, err := getRHACMConfigMapList(cmFlags)
	if err != nil || len(cms.Items) == 0 {
		return false
	}
	return true
}

func getRHACMConfigMapList(cmFlags *genericclioptions.CMFlags) (cms *corev1.ConfigMapList, err error) {
	f := cmFlags.KubectlFactory
	kubeClient, err := f.KubernetesClientSet()
	if err != nil {
		panic(err)
	}
	cms, err = kubeClient.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%v = %v", "ocm-configmap-type", "image-manifest"),
	})
	if (err != nil || len(cms.Items) == 0) && len(cmFlags.ServerNamespace) != 0 {
		cms, err = kubeClient.CoreV1().ConfigMaps(cmFlags.ServerNamespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%v = %v", "ocm-configmap-type", "image-manifest"),
		})
	}
	if err != nil || len(cms.Items) == 0 {
		cms, err = kubeClient.CoreV1().ConfigMaps("open-cluster-management").List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%v = %v", "ocm-configmap-type", "image-manifest"),
		})
	}
	return cms, err
}

func IsMCE(cmFlags *genericclioptions.CMFlags) bool {
	cms, err := getMCEConfigMapList(cmFlags)
	if err != nil || len(cms.Items) == 0 {
		return false
	}
	return true
}

func getMCEConfigMapList(cmFlags *genericclioptions.CMFlags) (cms *corev1.ConfigMapList, err error) {
	f := cmFlags.KubectlFactory
	kubeClient, err := f.KubernetesClientSet()
	if err != nil {
		panic(err)
	}
	cms, err = kubeClient.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%v = %v", "operators.coreos.com/multicluster-engine.multicluster-engine", ""),
	})
	if (err != nil || len(cms.Items) == 0) && len(cmFlags.ServerNamespace) != 0 {
		cms, err = kubeClient.CoreV1().ConfigMaps(cmFlags.ServerNamespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%v = %v", "operators.coreos.com/multicluster-engine.multicluster-engine", ""),
		})
	}
	if err != nil || len(cms.Items) == 0 {
		//TODO change default ns
		cms, err = kubeClient.CoreV1().ConfigMaps("multicluster-engine").List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%v = %v", "operators.coreos.com/multicluster-engine.multicluster-engine", ""),
		})
	}
	return cms, err
}
