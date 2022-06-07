// Copyright Contributors to the Open Cluster Management project
package managedcluster

import (
	"context"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"
	"github.com/stolostron/cm-cli/pkg/helpers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
)

func GetCmdAPIConfig(dynamicClient dynamic.Interface,
	kubeClient *kubernetes.Clientset,
	mc clusterv1.ManagedCluster,
	cph *clusterpoolhost.ClusterPoolHost) (config *clientcmdapi.Config, err error) {
	var cd *hivev1.ClusterDeployment
	var cc *hivev1.ClusterClaim
	foundOnCPH := false
	//Search clusterclaims in cph
	if cph != nil {
		cc, err = cph.GetClusterClaim(mc.Name, false, 60, false, nil)
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
		if cc != nil {
			foundOnCPH = true
		}
	}
	//Search for local clusterclaim
	if cc == nil {
		cc, err = getLocalClusterClaim(dynamicClient, mc)
		if err != nil {
			return nil, err
		}
	}
	if cc != nil {
		cd, err = cph.GetClusterDeployment(cc)
		if err != nil {
			return nil, err
		}
	}
	//Search for local clusterDeployment
	if cd == nil {
		cd, err = getLocalClusterDeployment(dynamicClient, mc)
		if err != nil {
			return nil, err
		}
	}
	if cd != nil {
		kubeConfigSecretName := cd.Spec.ClusterMetadata.AdminKubeconfigSecretRef.Name
		var kubeConfigSecret *corev1.Secret
		if foundOnCPH {
			clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
			if err != nil {
				return nil, err
			}
			kubeClientCPH, err := kubernetes.NewForConfig(clusterPoolRestConfig)
			if err != nil {
				return nil, err
			}
			kubeConfigSecret, err = kubeClientCPH.
				CoreV1().
				Secrets(cd.Namespace).
				Get(context.TODO(), kubeConfigSecretName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
		} else {
			kubeConfigSecret, err = kubeClient.
				CoreV1().
				Secrets(cd.Namespace).
				Get(context.TODO(), kubeConfigSecretName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
		}
		return clientcmd.Load(kubeConfigSecret.Data["kubeconfig"])
	}
	return nil, nil
}

func getLocalClusterClaim(dynamicClient dynamic.Interface, mc clusterv1.ManagedCluster) (*hivev1.ClusterClaim, error) {
	ccus, err := dynamicClient.Resource(helpers.GvrCC).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var ccu unstructured.Unstructured
	for _, ccui := range ccus.Items {
		if ccui.GetName() == mc.Name {
			ccu = ccui
		}
	}
	if ccu.Object != nil {
		cc := &hivev1.ClusterClaim{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc); err != nil {
			return nil, err
		}
		return cc, nil
	}
	return nil, nil
}

func getLocalClusterDeployment(dynamicClient dynamic.Interface, mc clusterv1.ManagedCluster) (*hivev1.ClusterDeployment, error) {
	cdus, err := dynamicClient.Resource(helpers.GvrCD).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var cdu unstructured.Unstructured
	for _, cdui := range cdus.Items {
		if cdui.GetName() == mc.Name {
			cdu = cdui
			break
		}
	}
	if cdu.Object != nil {
		cd := &hivev1.ClusterDeployment{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd); err != nil {
			return nil, err
		}
		return cd, nil
	}
	return nil, nil
}
