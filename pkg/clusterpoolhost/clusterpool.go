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
	pullSecret, err := kubeClient.CoreV1().Secrets("openshift-config").Get(
		context.TODO(),
		"pull-secret",
		metav1.GetOptions{})
	if err != nil {
		return err
	}

	ps, err := yaml.Marshal(pullSecret)
	if err != nil {
		return err
	}

	valueps := make(map[string]interface{})
	err = yaml.Unmarshal(ps, &valueps)
	if err != nil {
		return err
	}

	values["pullSecret"] = valueps

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
		"create/clusterpool/common/clusterimageset_cr.yaml",
		"create/clusterpool/common/clusterpool_cr.yaml",
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

func GetClusterPools(showCphName, dryRun bool) error {
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

	l, err := dynamicClient.Resource(helpers.GvrCP).Namespace(cph.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	if showCphName {
		fmt.Printf("%-20s\t%-32s\t%-4s\t%-5s\t%-12s\n", "CLUSTER_POOL_HOST", "CLUSTER_POOL", "SIZE", "READY", "ACTUAL_SIZE")
	} else {
		fmt.Printf("%-32s\t%-4s\t%-5s\t%-12s\n", "CLUSTER_POOL", "SIZE", "READY", "ACTUAL_SIZE")
	}
	if len(l.Items) == 0 {
		fmt.Printf("No clusterpool found for clusterpoolhost %s\n", cph.Name)
	}
	for _, cpu := range l.Items {
		cp := &hivev1.ClusterPool{}
		if runtime.DefaultUnstructuredConverter.FromUnstructured(cpu.UnstructuredContent(), cp); err != nil {
			return err
		}
		if showCphName {
			fmt.Printf("%-20s\t%-32s\t%-4d\t%-5d\t%-12d\n", cph.Name, cp.GetName(), cp.Spec.Size, cp.Status.Ready, cp.Status.Size)
		} else {
			fmt.Printf("%-32s\t%-4d\t%-5d\t%-12d\n", cp.GetName(), cp.Spec.Size, cp.Status.Ready, cp.Status.Size)
		}
	}
	return nil
}
