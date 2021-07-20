// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"context"
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
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
