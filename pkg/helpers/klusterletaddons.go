// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"context"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
	workclientset "open-cluster-management.io/api/client/work/clientset/versioned"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	workv1 "open-cluster-management.io/api/work/v1"
)

func WaitKlusterlet(clusterClient clusterclientset.Interface,
	clusterName string,
	timeout int) error {
	return wait.PollImmediate(10*time.Second, time.Duration(timeout)*time.Second, func() (bool, error) {
		return checkManagedClusterAvailable(clusterClient, clusterName)
	})
}

func WaitKlusterletAddons(workClient workclientset.Interface,
	clusterName string,
	timeout int) error {
	return wait.PollImmediate(10*time.Second, time.Duration(timeout)*time.Second, func() (bool, error) {
		return checkManagedClusterAddons(workClient, clusterName)
	})
}

func checkManagedClusterAvailable(clusterClient clusterclientset.Interface, clusterName string) (bool, error) {
	mc, err := clusterClient.ClusterV1().ManagedClusters().Get(context.TODO(), clusterName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	for _, condition := range mc.Status.Conditions {
		if condition.Type == clusterv1.ManagedClusterConditionAvailable &&
			condition.Status == metav1.ConditionTrue {
			fmt.Printf("agent on %s ready\n", clusterName)
			return true, nil
		}
	}
	fmt.Printf("agent on %s not ready\n", clusterName)
	return false, nil
}

func checkManagedClusterAddons(workClient workclientset.Interface, clusterName string) (bool, error) {
	mwl, err := workClient.WorkV1().ManifestWorks(clusterName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	allReady := true
	for _, mw := range mwl.Items {
		if !strings.HasPrefix(mw.Name, clusterName+"-klusterlet-addon-") {
			continue
		}
		addon := strings.TrimPrefix(mw.Name, clusterName+"-klusterlet-addon-")
		addonReady := false
		for _, condition := range mw.Status.Conditions {
			if condition.Type == workv1.WorkAvailable &&
				condition.Status == metav1.ConditionTrue {
				addonReady = true
				fmt.Printf("addon %s on %s ready\n", addon, clusterName)
			}
			fmt.Printf("addon %s on %s not ready\n", addon, clusterName)
		}
		if !addonReady {
			allReady = false
		}
	}
	return allReady, nil
}
