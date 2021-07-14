// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost/scenario"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"
)

var (
	gvrCC schema.GroupVersionResource = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterclaims"}
	gvrCD schema.GroupVersionResource = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterdeployments"}
)

func CreateClusterClaims(clusterClaimNames, clusterPoolName string, skipSchedule bool, timeout int, dryRun bool, outputFile string) error {
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

	_, err = dynamicClient.Resource(gvrCP).Namespace(cph.Namespace).Get(context.TODO(), clusterPoolName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	me, err := WhoAmI(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	serviceAccountName := strings.TrimPrefix(me.Name, "system:serviceaccount:"+cph.Namespace+":")

	reader := scenario.GetScenarioResourcesReader()

	output := make([]string, 0)
	for _, ccn := range strings.Split(clusterClaimNames, ",") {
		clusterClaimName := strings.TrimSpace(ccn)
		values := make(map[string]string)
		values["Name"] = clusterClaimName
		values["Namespace"] = cph.Namespace
		values["ClusterPoolName"] = clusterPoolName
		values["ServiceAccountName"] = serviceAccountName
		values["Group"] = cph.Group
		files := []string{
			"create/clusterclaim/clusterclaim_cr.yaml",
		}

		applierBuilder := &clusteradmapply.ApplierBuilder{}
		applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient)
		out, err := applier.ApplyCustomResources(reader, values, dryRun, "", files...)
		if err != nil {
			return err
		}
		fmt.Printf("clusterclaim %s created\n", clusterClaimName)
		output = append(output, out...)
	}

	if !dryRun {
		if err := waitClusterClaimsRunning(dynamicClient, clusterClaimNames, clusterPoolName, cph.Namespace, timeout); err != nil {
			return err
		}
	}
	return clusteradmapply.WriteOutput(outputFile, output)
}

func RunClusterClaims(clusterClaimNames string, skipSchedule bool, timeout int, dryRun bool, outputFile string) error {
	skipScheduleAction := "true"
	if skipSchedule {
		skipScheduleAction = "skip"
	}
	cph, err := GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}
	if err = cph.setHibernateClusterClaims(clusterClaimNames, false, skipScheduleAction, dryRun, outputFile); err != nil {
		return err
	}
	if !dryRun {
		clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
		if err != nil {
			return err
		}

		dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
		if err != nil {
			return err
		}
		return waitClusterClaimsRunning(dynamicClient, clusterClaimNames, "", cph.Namespace, timeout)
	}
	return nil
}

func HibernateClusterClaims(clusterClaimNames string, skipSchedule, dryRun bool, outputFile string) error {
	skipScheduleAction := "true"
	if skipSchedule {
		skipScheduleAction = "skip"
	}
	cph, err := GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}
	return cph.setHibernateClusterClaims(clusterClaimNames, true, skipScheduleAction, dryRun, outputFile)
}

func (cph *ClusterPoolHost) setHibernateClusterClaims(clusterClaimNames string, hibernate bool, skipScheduleAction string, dryRun bool, outputFile string) error {
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	for _, ccn := range strings.Split(clusterClaimNames, ",") {
		ccu, err := dynamicClient.Resource(gvrCC).Namespace(cph.Namespace).Get(context.TODO(), ccn, metav1.GetOptions{})
		if err != nil {
			return err
		}
		cc := &hivev1.ClusterClaim{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc); err != nil {
			return err
		}
		cdu, err := dynamicClient.Resource(gvrCD).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if !dryRun {
			cd := &hivev1.ClusterDeployment{}
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd)
			if err != nil {
				return err
			}
			if hibernate {
				cd.Spec.PowerState = hivev1.HibernatingClusterPowerState
			} else {
				cd.Spec.PowerState = hivev1.RunningClusterPowerState
			}
			if len(skipScheduleAction) != 0 {
				cd.Labels["hibernate"] = skipScheduleAction
			}
			cdu.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(cd)
			if err != nil {
				return err
			}
			_, err = dynamicClient.Resource(gvrCD).Namespace(cc.Spec.Namespace).Update(context.TODO(), cdu, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func waitClusterClaimsRunning(dynamicClient dynamic.Interface, clusterClaimNames, clusterPoolName, namespace string, timeout int) error {
	i := 0
	return wait.PollImmediate(1*time.Minute, time.Duration(timeout)*time.Minute, func() (bool, error) {
		i += 1
		return checkClusterClaimsRunning(dynamicClient, clusterClaimNames, clusterPoolName, namespace, i, timeout)
	})

}
func checkClusterClaimsRunning(dynamicClient dynamic.Interface, clusterClaimNames, clusterPoolName, namespace string, i, timeout int) (bool, error) {
	if len(clusterPoolName) != 0 {
		cpu, err := dynamicClient.Resource(gvrCP).Namespace(namespace).Get(context.TODO(), clusterPoolName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		cp := &hivev1.ClusterPool{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(cpu.UnstructuredContent(), cp)
		if err != nil {
			return false, err
		}
		if cp.Spec.Size == 0 {
			fmt.Printf("WARNING: the clusterpool %s size is 0, should be at least 1 for the clusterclaim to be honored\n", clusterPoolName)
		}
	}
	allErrors := make(map[string]error)
	allRunning := true
	for _, ccn := range strings.Split(clusterClaimNames, ",") {
		clusterClaimName := strings.TrimSpace(ccn)
		ccu, err := dynamicClient.Resource(gvrCC).Namespace(namespace).Get(context.TODO(), clusterClaimName, metav1.GetOptions{})
		if err != nil {
			allErrors[clusterClaimName] = err
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}
		cc := &hivev1.ClusterClaim{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc)
		if err != nil {
			allErrors[clusterClaimName] = err
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}
		running := false
		if len(cc.Spec.Namespace) != 0 {
			cdu, err := dynamicClient.Resource(gvrCD).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
			if err != nil {
				allErrors[clusterClaimName] = err
				fmt.Printf("Error: %s\n", err.Error())
				continue
			}
			cd := &hivev1.ClusterDeployment{}
			if runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd); err != nil {
				allErrors[clusterClaimName] = err
				fmt.Printf("Error: %s\n", err.Error())
				continue
			}
			if cd.Spec.PowerState == hivev1.HibernatingClusterPowerState {
				allErrors[clusterClaimName] = fmt.Errorf("%s is hibernating, run a use command to resume it", cc.GetName())
				fmt.Printf("Error: %s\n", allErrors[clusterClaimName])
				continue
			}
			c := getClusterClaimRunningStatus(cc)
			if len(cd.Spec.ClusterMetadata.AdminPasswordSecretRef.Name) != 0 &&
				len(cd.Spec.BaseDomain) != 0 &&
				len(cd.Status.APIURL) != 0 &&
				c != nil && c.Status == corev1.ConditionStatus(metav1.ConditionTrue) {
				running = true
				fmt.Printf("clusterclaim %s is running with id %s (%d/%d)\n", clusterClaimName, cc.Spec.Namespace, i, timeout)
			}
		}
		if !running {
			fmt.Printf("clusterclaim %s is not running (%d/%d)\n", clusterClaimName, i, timeout)
			allRunning = false
		}
	}
	if len(allErrors) == len(strings.Split(clusterClaimNames, ",")) {
		return false, fmt.Errorf("all requested clusterclaims have errors")
	}
	if len(allErrors) == 0 {
		return allRunning, nil
	}
	return false, nil
}

func getClusterClaimRunningStatus(cc *hivev1.ClusterClaim) *hivev1.ClusterClaimCondition {
	for _, c := range cc.Status.Conditions {
		if c.Type == hivev1.ClusterRunningCondition {
			return &c
		}
	}
	return nil
}

func DeleteClusterClaims(clusterClaimNames string, dryRun bool, outputFile string) error {
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
		gvr := schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterclaims"}
		for _, ccn := range strings.Split(clusterClaimNames, ",") {
			clusterClaimName := strings.TrimSpace(ccn)
			err = dynamicClient.Resource(gvr).Namespace(cph.Namespace).Delete(context.TODO(), clusterClaimName, metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetClusterClaims(showCphName, dryRun bool, outputFile string) error {
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

	l, err := dynamicClient.Resource(gvrCC).Namespace(cph.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	if showCphName {
		fmt.Printf("%-20s\t%-15s\t%-11s\t%-9s\t%-20s\n", "CLUSTER_POOL_HOST", "CLUSTER_CLAIM", "POWER_STATE", "HIBERNATE", "ID")
	} else {
		fmt.Printf("%-15s\t%-11s\t%-9s\t%-20s\n", "CLUSTER_CLAIM", "POWER_STATE", "HIBERNATE", "ID")
	}
	if len(l.Items) == 0 {
		fmt.Printf("No clusterclaim found for clusterpoolhost %s\n", cph.Name)
	}
	for _, ccu := range l.Items {
		cc := &hivev1.ClusterClaim{}
		if runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc); err != nil {
			return err
		}
		c := getClusterClaimPendingStatus(cc)
		if c != nil && c.Status == corev1.ConditionStatus(metav1.ConditionTrue) {
			printClusterClaim(showCphName, true, cph, cc, nil)
			continue
		}
		cdu, err := dynamicClient.Resource(gvrCD).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
		if err != nil {
			if showCphName {
				fmt.Printf("%s clusterdeployment %s\n", cc.GetName(), err.Error())
			} else {
				fmt.Printf("%s %s clusterdeployment %s\n", cph.Name, cc.GetName(), err.Error())
			}
			continue
		}
		cd := &hivev1.ClusterDeployment{}
		if runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd); err != nil {
			return err
		}
		printClusterClaim(showCphName, false, cph, cc, cd)
	}
	return nil
}

func getClusterClaimPendingStatus(cc *hivev1.ClusterClaim) *hivev1.ClusterClaimCondition {
	for _, c := range cc.Status.Conditions {
		if c.Type == hivev1.ClusterClaimPendingCondition {
			return &c
		}
	}
	return nil
}

func printClusterClaim(showCphName bool, pending bool, cph *ClusterPoolHost, cc *hivev1.ClusterClaim, cd *hivev1.ClusterDeployment) {
	var powerState, hibernate, cdName string
	if cd != nil {
		powerState = string(cd.Spec.PowerState)
		hibernate = cd.Labels["hibernate"]
		cdName = cd.Name
	}
	if pending {
		powerState = string(hivev1.ClusterClaimPendingCondition)
	}
	if showCphName {
		fmt.Printf("%-20s\t%-15s\t%-11s\t%-9s\t%-20s\n", cph.Name, cc.GetName(), powerState, hibernate, cdName)
	} else {
		fmt.Printf("%-15s\t%-11s\t%-9s\t%-20s\n", cc.GetName(), powerState, hibernate, cdName)
	}
}

func GetClusterClaim(clusterName string, timeout int, dryRun bool) error {
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
	klog.V(3).Infof("Wait cc %s ready", clusterName)
	if err = waitClusterClaimsRunning(dynamicClient, clusterName, "", cph.Namespace, timeout); err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}
	ccu, err := dynamicClient.Resource(gvrCC).Namespace(cph.Namespace).Get(context.TODO(), clusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cc := &hivev1.ClusterClaim{}
	if runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc); err != nil {
		return err
	}
	cdu, err := dynamicClient.Resource(gvrCD).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cd := &hivev1.ClusterDeployment{}
	if runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd); err != nil {
		return err
	}
	if cd.Spec.PowerState == hivev1.HibernatingClusterPowerState {
		return fmt.Errorf("%s is hibernating, run a use command to resume it", cc.GetName())
	}
	s, err := kubeClient.CoreV1().Secrets(cd.Namespace).Get(context.TODO(), cd.Spec.ClusterMetadata.AdminPasswordSecretRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("username:    %s\n", s.Data["username"])
	fmt.Printf("password:    %s\n", s.Data["password"])
	fmt.Printf("basedomain:  %s\n", cd.Spec.BaseDomain)
	fmt.Printf("api_url:     %s\n", cd.Status.APIURL)
	fmt.Printf("console_url: %s\n", cd.Status.WebConsoleURL)
	return nil
}
