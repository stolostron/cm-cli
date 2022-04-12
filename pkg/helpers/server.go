// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"context"
	"fmt"

	"github.com/stolostron/cm-cli/pkg/genericclioptions"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"
)

func IsSupported(cmFlags *genericclioptions.CMFlags) (isSupported bool, err error) {
	return IsRHACM(cmFlags) || IsMCE(cmFlags), err
}

func IsRHACM(cmFlags *genericclioptions.CMFlags) bool {
	if cmFlags.SkipServerCheck {
		return true
	}
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
	if cmFlags.SkipServerCheck {
		return true
	}
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
		LabelSelector: fmt.Sprintf("%v = %v", "operators.coreos.com/.multicluster-engine", ""),
	})
	if (err != nil || len(cms.Items) == 0) && len(cmFlags.ServerNamespace) != 0 {
		cms, err = kubeClient.CoreV1().ConfigMaps(cmFlags.ServerNamespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%v = %v", "operators.coreos.com/multicluster-engine."+cmFlags.ServerNamespace, ""),
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

func IsOpenshift(cmFlags *genericclioptions.CMFlags) (bool, error) {
	_, _, dynamicClient, err := clusteradmhelpers.GetClients(cmFlags.KubectlFactory)
	if err != nil {
		return false, err
	}
	_, err = dynamicClient.Resource(GvrOpenshiftClusterVersions).Get(context.TODO(), "version", metav1.GetOptions{})
	return err == nil, nil
}

func IsHypershift(cmFlags *genericclioptions.CMFlags) (bool, error) {
	_, apiExtensionClient, _, err := clusteradmhelpers.GetClients(cmFlags.KubectlFactory)
	if err != nil {
		return false, err
	}

	_, err = apiExtensionClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), GvrHC.GroupResource().String(), metav1.GetOptions{})
	return err == nil, nil
}
