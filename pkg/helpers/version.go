// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"

	"k8s.io/client-go/kubernetes"
)

func GetACMVersion(kubeClient kubernetes.Interface, dynamicClient dynamic.Interface) (version, snapshot string, err error) {
	lo := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%v = %v", "ocm-configmap-type", "image-manifest"),
	}
	cms, err := kubeClient.CoreV1().ConfigMaps("").List(context.TODO(), lo)
	if err != nil {
		return "", "", err
	}
	if len(cms.Items) == 0 {
		return "", "", fmt.Errorf("no configmap with labelset %v", lo.LabelSelector)
	}
	ns := cms.Items[0].Namespace

	umch, err := dynamicClient.Resource(GvrMCH).Namespace(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", "", err
	}
	if len(umch.Items) == 0 {
		return "", "", fmt.Errorf("no multiclusterhub found in namespace %s", ns)
	}
	ustatus, ok := umch.Items[0].Object["status"]
	if !ok {
		return "", "", fmt.Errorf("no status found multiclusterhub in %s/%s", ns, umch.Items[0].GetName())
	}
	uversion, ok := ustatus.(map[string]interface{})["currentVersion"]
	if !ok {
		return "", "", fmt.Errorf("no currentVersion found multiclusterhub in %s/%s", ns, umch.Items[0].GetName())
	}
	version = uversion.(string)
	lo = metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%v = %v,%v = %v", "ocm-configmap-type", "image-manifest", "ocm-release-version", version),
	}
	cms, err = kubeClient.CoreV1().ConfigMaps("").List(context.TODO(), lo)
	if err != nil {
		return version, "", err
	}
	if len(cms.Items) == 1 {
		ns := cms.Items[0].Namespace
		if v, ok := cms.Items[0].Labels["ocm-release-version"]; ok {
			version = v
		}
		acmRegistryDeployment, err := kubeClient.AppsV1().Deployments(ns).Get(context.TODO(), "acm-custom-registry", metav1.GetOptions{})
		if err == nil {
			for _, c := range acmRegistryDeployment.Spec.Template.Spec.Containers {
				if strings.Contains(c.Image, "acm-custom-registry") {
					snapshot = strings.Split(c.Image, ":")[1]
					break
				}
			}
		}
	}
	return version, snapshot, nil
}

func IsSupported(kubeClient kubernetes.Interface, dynamicClient dynamic.Interface, constraint string) (bool, error) {
	version, _, err := GetACMVersion(kubeClient, dynamicClient)
	if err != nil {
		return false, err
	}

	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return false, err
	}

	vs, err := semver.NewVersion(version)
	if err != nil {
		return false, err
	}

	return c.Check(vs), nil
}
