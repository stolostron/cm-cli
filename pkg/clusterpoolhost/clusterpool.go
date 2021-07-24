// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"
)

func SizeClusterPool(clusterPoolName string, size int32, dryRun bool) error {
	cph, err := GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}
	cpu, err := dynamicClient.Resource(helpers.GvrCP).Namespace(cph.Namespace).Get(context.TODO(), clusterPoolName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cp := &hivev1.ClusterPool{}
	if runtime.DefaultUnstructuredConverter.FromUnstructured(cpu.UnstructuredContent(), cp); err != nil {
		return err
	}
	if !dryRun {
		cp.Spec.Size = size
		cpu.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(cp)
		if runtime.DefaultUnstructuredConverter.FromUnstructured(cpu.UnstructuredContent(), cp); err != nil {
			return err
		}
		if _, err = dynamicClient.Resource(helpers.GvrCP).Namespace(cph.Namespace).Update(context.TODO(), cpu, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func CreateClusterPool(clusterPoolName, cloud string, values map[string]interface{}, dryRun bool, outputFile string) error {
	cph, err := GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}
	values["namespace"] = cph.Namespace

	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	apiExtensionsClient, err := apiextensionsclient.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	output := make([]string, 0)

	reader := scenario.GetScenarioResourcesReader()
	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient).Build()

	installConfig, err := applier.MustTempalteAsset(reader, values, "", filepath.Join("create", "clusterpool", cloud, "install_config.yaml"))
	if err != nil {
		return err
	}

	valueic := make(map[string]interface{})
	err = yaml.Unmarshal(installConfig, &valueic)
	if err != nil {
		return err
	}

	files := []string{
		"create/clusterpool/common/namespace.yaml",
	}

	out, err := applier.ApplyDirectly(reader, values, dryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)
	values["installConfig"] = valueic

	files = []string{
		"create/clusterpool/common/creds_secret_cr.yaml",
		"create/clusterpool/common/install_config_secret_cr.yaml",
		"create/clusterpool/common/pull_secret_cr.yaml",
	}

	cpi := values["clusterPool"]
	cp := cpi.(map[string]interface{})
	if _, ok := cp["imageSetRef"]; ok {
		files = append(files,
			"create/clusterpool/common/clusterpool_cr.yaml")
	} else {
		files = append(files, "create/clusterpool/common/clusterimageset_cr.yaml",
			"create/clusterpool/common/clusterpool_cr.yaml")
	}
	out, err = applier.ApplyCustomResources(reader, values, dryRun, "create/clusterpool/common/_helpers.tpl", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	return clusteradmapply.WriteOutput(outputFile, output)
}

func DeleteClusterPools(clusterPoolNames string, dryRun bool, outputFile string) error {
	cph, err := GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	if !dryRun {
		for _, cpn := range strings.Split(clusterPoolNames, ",") {
			clusterPoolName := strings.TrimSpace(cpn)
			err = dynamicClient.Resource(helpers.GvrCP).Namespace(cph.Namespace).Delete(context.TODO(), clusterPoolName, metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetClusterPools(showCphName, dryRun bool) (*hivev1.ClusterPoolList, error) {
	clusterPools := &hivev1.ClusterPoolList{}
	cph, err := GetCurrentClusterPoolHost()
	if err != nil {
		return clusterPools, err
	}
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return clusterPools, err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return clusterPools, err
	}

	l, err := dynamicClient.Resource(helpers.GvrCP).Namespace(cph.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return clusterPools, err
	}
	for _, cpu := range l.Items {
		cp := &hivev1.ClusterPool{}
		if runtime.DefaultUnstructuredConverter.FromUnstructured(cpu.UnstructuredContent(), cp); err != nil {
			return clusterPools, err
		}
		clusterPools.Items = append(clusterPools.Items, *cp)
	}
	return clusterPools, nil
}

func SprintClusterPools(cph *ClusterPoolHost, sep string, cps *hivev1.ClusterPoolList) []string {
	lines := make([]string, 0)
	if len(cps.Items) != 0 {
		lines = append(lines, fmt.Sprintf("%s%s%s%s%s%s%s%s%s", "CLUSTER_POOL_HOST", sep, "CLUSTER_POOL", sep, "SIZE", sep, "READY", sep, "ACTUAL_SIZE"))
	}
	for _, cp := range cps.Items {
		lines = append(lines, sprintClusterPool(cph, sep, &cp))
	}
	if len(cps.Items) != 0 {
		lines = append(lines, fmt.Sprintf("%s%s%s%s%s%s%s%s%s", "", sep, "", sep, "", sep, "", sep, ""))
	}
	klog.V(5).Infof("lines:%s\n", lines)
	return lines
}

func sprintClusterPool(cph *ClusterPoolHost, sep string, cp *hivev1.ClusterPool) string {
	return fmt.Sprintf("%s%s%s%s%4d%s%5d%s%11d", cph.Name, sep, cp.GetName(), sep, cp.Spec.Size, sep, cp.Status.Ready, sep, cp.Status.Size)
}
