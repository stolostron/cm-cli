// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

func GetACMVersion(kubeClient kubernetes.Interface) (version, snapshot string, err error) {
	lo := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%v = %v", "ocm-configmap-type", "image-manifest"),
	}
	cms, err := kubeClient.CoreV1().ConfigMaps("").List(context.TODO(), lo)
	if err != nil {
		return "", "", err
	}
	if len(cms.Items) > 1 {
		return "", "", fmt.Errorf("found more than one configmap with labelset %v", lo.LabelSelector)
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

func IsSupported(kubeClient kubernetes.Interface, constraint string) (bool, error) {
	version, _, err := GetACMVersion(kubeClient)
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
