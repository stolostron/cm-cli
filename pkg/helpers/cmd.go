// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func GetExampleHeader() string {
	switch os.Args[0] {
	case "oc":
		return "oc cm"
	case "kubectl":
		return "kubectl cm"
	default:
		return os.Args[0]
	}
}

func IsRHACM(f cmdutil.Factory) bool {
	kubeClient, err := f.KubernetesClientSet()
	if err != nil {
		panic(err)
	}
	cms, err := kubeClient.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%v = %v", "ocm-configmap-type", "image-manifest"),
	})
	if err != nil {
		return false
	}
	if len(cms.Items) == 0 {
		return false
	}
	return true
}

func IsMCE(f cmdutil.Factory) bool {
	kubeClient, err := f.KubernetesClientSet()
	if err != nil {
		panic(err)
	}
	cms, err := kubeClient.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%v = %v", "operators.coreos.com/multicluster-engine.multicluster-engine", ""),
	})
	if err != nil {
		return false
	}
	if len(cms.Items) == 0 {
		return false
	}
	return true
}
