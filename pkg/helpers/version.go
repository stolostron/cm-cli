// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/stolostron/cm-cli/pkg/genericclioptions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"

	"k8s.io/client-go/kubernetes"
)

const (
	RHACM      string = "RHACM"
	MCE        string = "MCE"
	HYPERSHIFT string = "Hypershift"
)

func GetACMVersion(cmFlags *genericclioptions.CMFlags, kubeClient kubernetes.Interface, dynamicClient dynamic.Interface) (version, snapshot string, err error) {
	cms, err := getRHACMConfigMapList(cmFlags)
	if err != nil {
		return "", "", err
	}
	if len(cms.Items) == 0 {
		return "", "", fmt.Errorf("no configmap found")
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
	lo := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%v = %v,%v = %v", "ocm-configmap-type", "image-manifest", "ocm-release-version", version),
	}
	cms, err = kubeClient.CoreV1().ConfigMaps(ns).List(context.TODO(), lo)
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

func GetMCEVersion(cmFlags *genericclioptions.CMFlags, kubeClient kubernetes.Interface, dynamicClient dynamic.Interface) (version, snapshot string, err error) {
	cms, err := getMCEConfigMapList(cmFlags)
	if err != nil {
		return "", "", err
	}
	if len(cms.Items) == 0 {
		return "", "", fmt.Errorf("no configmap found")
	}
	ns := cms.Items[0].Namespace
	ucsv, err := dynamicClient.Resource(GvrCSV).Namespace(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", "", err
	}
	if err != nil {
		return "", "", err
	}
	if len(ucsv.Items) == 0 {
		return "", "", fmt.Errorf("no clusterserviceversion found in namespace %s", ns)
	}

	uspec := ucsv.Items[0].Object["spec"]
	spec := uspec.(map[string]interface{})
	version = spec["version"].(string)
	return version, "", nil

}

func GetVersion(cmFlags *genericclioptions.CMFlags, isCPHCommand bool, cphName string) (version string, platform string, err error) {
	f := cmFlags.KubectlFactory
	kubeClient, err := f.KubernetesClientSet()
	if err != nil {
		return version, platform, err
	}
	dynamicClient, err := f.DynamicClient()
	if err != nil {
		return version, platform, err
	}
	switch {
	case IsRHACM(cmFlags):
		platform = RHACM
		version, _, err = GetACMVersion(cmFlags, kubeClient, dynamicClient)
	case IsMCE(cmFlags):
		platform = MCE
		version, _, err = GetMCEVersion(cmFlags, kubeClient, dynamicClient)
	}
	return version, platform, err
}

func IsSupportedVersion(cmFlags *genericclioptions.CMFlags, isCPHCommand bool, cphName string, rhacmConstraint string, mceConstraint string) (isSupported bool, platform string, err error) {
	var version string
	f := cmFlags.KubectlFactory
	kubeClient, err := f.KubernetesClientSet()
	if err != nil {
		return isSupported, platform, err
	}
	dynamicClient, err := f.DynamicClient()
	if err != nil {
		return isSupported, platform, err
	}
	var c *semver.Constraints
	switch {
	case IsRHACM(cmFlags):
		platform = RHACM
		if len(rhacmConstraint) == 0 {
			return false, platform, nil
		}
		version, _, err = GetACMVersion(cmFlags, kubeClient, dynamicClient)
		if err != nil {
			return isSupported, platform, err
		}
		c, err = semver.NewConstraint(rhacmConstraint)
	case IsMCE(cmFlags):
		platform = MCE
		if len(mceConstraint) == 0 {
			return false, platform, nil
		}
		version, _, err = GetMCEVersion(cmFlags, kubeClient, dynamicClient)
		if err != nil {
			return isSupported, platform, err
		}
		c, err = semver.NewConstraint(mceConstraint)
	}
	if err != nil {
		return isSupported, platform, err
	}

	vs, err := semver.NewVersion(version)
	if err != nil {
		return isSupported, platform, err
	}

	return c.Check(vs), platform, nil
}
