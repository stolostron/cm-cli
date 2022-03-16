// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	printclusterpoolv1alpha1 "github.com/stolostron/cm-cli/api/cm-cli/v1alpha1"
	"github.com/stolostron/cm-cli/pkg/clusterpoolhost/scenario"
	"github.com/stolostron/cm-cli/pkg/helpers"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/cmd/get"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"
)

func (cph *ClusterPoolHost) GetClusterContextName(clusterName string) string {
	return fmt.Sprintf("%s/%s", cph.Name, clusterName)
}

func (cph *ClusterPoolHost) CreateClusterClaims(clusterClaimNames, clusterPoolName string, skipSchedule bool, autoImport bool, timeout int, dryRun bool, outputFile string, printFlags *get.PrintFlags) error {
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

	_, err = dynamicClient.Resource(helpers.GvrCP).Namespace(cph.Namespace).Get(context.TODO(), clusterPoolName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	//Check if all clusterclaims are not present before creation
	for _, ccn := range strings.Split(clusterClaimNames, ",") {
		clusterClaimName := strings.TrimSpace(ccn)
		cc, err := dynamicClient.Resource(helpers.GvrCC).Namespace(cph.Namespace).Get(context.TODO(), clusterClaimName, metav1.GetOptions{})
		switch {
		case errors.IsNotFound(err):
		case err == nil:
			if cc.GetDeletionTimestamp().IsZero() {
				return fmt.Errorf("clusterclaim %s already exists", ccn)
			} else {
				return fmt.Errorf("clusterclaim %s is currently in deletion, please retry later", ccn)
			}
		case err != nil:
			return err
		}
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
		values["AutoImport"] = strconv.FormatBool(autoImport)
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
		if err := waitClusterClaimsRunning(dynamicClient, clusterClaimNames, clusterPoolName, cph.Namespace, timeout, printFlags); err != nil {
			return err
		}
	}
	return clusteradmapply.WriteOutput(outputFile, output)
}

func (cph *ClusterPoolHost) RunClusterClaims(clusterClaimNames string, skipSchedule bool, timeout int, dryRun bool, outputFile string, printFlags *get.PrintFlags) error {
	skipScheduleAction := "true"
	if skipSchedule {
		skipScheduleAction = "skip"
	}
	if err := cph.setHibernateClusterClaims(clusterClaimNames, false, skipScheduleAction, dryRun, outputFile); err != nil {
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
		return waitClusterClaimsRunning(dynamicClient, clusterClaimNames, "", cph.Namespace, timeout, printFlags)
	}
	return nil
}

func (cph *ClusterPoolHost) HibernateClusterClaims(clusterClaimNames string, skipSchedule, dryRun bool, outputFile string) error {
	skipScheduleAction := "true"
	if skipSchedule {
		skipScheduleAction = "skip"
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
		ccu, err := dynamicClient.Resource(helpers.GvrCC).Namespace(cph.Namespace).Get(context.TODO(), ccn, metav1.GetOptions{})
		if err != nil {
			return err
		}
		cc := &hivev1.ClusterClaim{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc); err != nil {
			return err
		}
		if len(cc.Spec.Namespace) == 0 {
			return fmt.Errorf("something wrong happened, the clusterclaim %s doesn't have a spec.namespace set", cc.Name)
		}
		cdu, err := dynamicClient.Resource(helpers.GvrCD).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
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
				cd.Spec.PowerState = hivev1.ClusterPowerStateHibernating
			} else {
				cd.Spec.PowerState = hivev1.ClusterPowerStateRunning
			}
			if len(skipScheduleAction) != 0 {
				cd.Labels["hibernate"] = skipScheduleAction
			}
			cdu.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(cd)
			if err != nil {
				return err
			}
			_, err = dynamicClient.Resource(helpers.GvrCD).Namespace(cc.Spec.Namespace).Update(context.TODO(), cdu, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func waitClusterClaimsRunning(dynamicClient dynamic.Interface, clusterClaimNames, clusterPoolName, namespace string, timeout int, printFlags *get.PrintFlags) error {
	if timeout == 0 {
		running, err := checkClusterClaimsRunning(dynamicClient, clusterClaimNames, clusterPoolName, namespace, 0, timeout, false, printFlags)
		if err != nil {
			return err
		}
		if !running {
			return fmt.Errorf("one or more clusterclaims are not running: %s", clusterClaimNames)
		}
		return nil
	}
	i := 0
	return wait.PollImmediate(1*time.Minute, time.Duration(timeout)*time.Minute, func() (bool, error) {
		running, err := checkClusterClaimsRunning(dynamicClient, clusterClaimNames, clusterPoolName, namespace, i, timeout, false, printFlags)
		i += 1
		return running, err
	})

}

func checkClusterClaimsRunning(dynamicClient dynamic.Interface, clusterClaimNames, clusterPoolName, namespace string, i, timeout int, errorOnHibernate bool, printFlags *get.PrintFlags) (bool, error) {
	if len(clusterPoolName) != 0 {
		cpu, err := dynamicClient.Resource(helpers.GvrCP).Namespace(namespace).Get(context.TODO(), clusterPoolName, metav1.GetOptions{})
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
		ccu, err := dynamicClient.Resource(helpers.GvrCC).Namespace(namespace).Get(context.TODO(), clusterClaimName, metav1.GetOptions{})
		if err != nil {
			allErrors[clusterClaimName] = fmt.Errorf("(%d/%d) clusterclaim %s error: %s", i, timeout, clusterClaimName, err.Error())
			continue
		}
		cc := &hivev1.ClusterClaim{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc)
		if err != nil {
			allErrors[clusterClaimName] = fmt.Errorf("(%d/%d) clusterclaim %s error: %s", i, timeout, clusterClaimName, err.Error())
			continue
		}
		running, allErrorsO, err := checkClusterClaimRunning(dynamicClient, clusterClaimName, namespace, cc, i, timeout, errorOnHibernate, printFlags)
		if err != nil {
			allErrors[clusterClaimName] = err
		}
		for k, v := range allErrorsO {
			allErrors[k] = v
		}
		if !running {
			if timeout == 0 {
				allErrors[clusterClaimName] = fmt.Errorf("(%d) clusterclaim %s is not running", i, clusterClaimName)
			} else {
				allErrors[clusterClaimName] = fmt.Errorf("(%d/%d) clusterclaim %s is not running", i, timeout, clusterClaimName)
			}
			allRunning = false
		} else {
			delete(allErrors, clusterClaimName)
			if printFlags == nil || printFlags.OutputFormat == nil || strings.HasPrefix(*printFlags.OutputFormat, "custom-columns=") {
				if timeout == 0 {
					fmt.Printf("(%d) clusterclaim %s is running with id %s\n", i, clusterClaimName, cc.Spec.Namespace)
				} else {
					fmt.Printf("(%d/%d) clusterclaim %s is running with id %s\n", i, timeout, clusterClaimName, cc.Spec.Namespace)
				}
			}
		}
	}
	for _, msg := range allErrors {
		fmt.Println(msg)
	}
	if len(allErrors) == 0 {
		return allRunning, nil
	}
	return false, nil
}

func checkClusterClaimRunning(dynamicClient dynamic.Interface, clusterClaimName, namespace string, cc *hivev1.ClusterClaim, i, timeout int, errorOnHibernate bool, printFlags *get.PrintFlags) (running bool, allErrors map[string]error, err error) {
	running = false
	allErrors = make(map[string]error)
	if len(cc.Spec.Namespace) != 0 {
		cdu, err := dynamicClient.Resource(helpers.GvrCD).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
		if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonForbidden {
			allErrors[clusterClaimName] = fmt.Errorf("permissions error when accessing claimed ClusterDeployment.  permissions are likely still propagating. \nerror: %s", err.Error())
			return running, allErrors, nil
		}
		if err != nil {
			return running, allErrors, fmt.Errorf("(%d/%d) clusterclaim %s error: %s", i, timeout, clusterClaimName, err.Error())
		}
		cd := &hivev1.ClusterDeployment{}
		if runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd); err != nil {
			return running, allErrors, fmt.Errorf("(%d/%d) clusterclaim %s error: %s", i, timeout, clusterClaimName, err.Error())
		}
		if errorOnHibernate && cd.Spec.PowerState == hivev1.ClusterPowerStateHibernating {
			return running, allErrors, fmt.Errorf("(%d/%d) clusterclaim %s is hibernating, run a \"cm use cc\" or \"cm run cc\" command to resume it", i, timeout, cc.GetName())
		}
		c := getClusterClaimRunningStatus(cc)
		if len(cd.Spec.ClusterMetadata.AdminPasswordSecretRef.Name) != 0 &&
			len(cd.Spec.BaseDomain) != 0 &&
			len(cd.Status.APIURL) != 0 &&
			c != nil && c.Status == corev1.ConditionStatus(metav1.ConditionTrue) {
			running = true
		}
	}
	return running, allErrors, nil
}

func getClusterClaimRunningStatus(cc *hivev1.ClusterClaim) *hivev1.ClusterClaimCondition {
	for _, c := range cc.Status.Conditions {
		if c.Type == hivev1.ClusterRunningCondition {
			return &c
		}
	}
	return nil
}

func (cph *ClusterPoolHost) DeleteClusterClaims(clusterClaimNames string, dryRun bool, outputFile string) error {
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	if !dryRun {
		for _, ccn := range strings.Split(clusterClaimNames, ",") {
			clusterClaimName := strings.TrimSpace(ccn)
			err = dynamicClient.Resource(helpers.GvrCC).Namespace(cph.Namespace).Delete(context.TODO(), clusterClaimName, metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (cph *ClusterPoolHost) GetClusterClaims(dryRun bool) (*hivev1.ClusterClaimList, error) {
	clusterClaims := &hivev1.ClusterClaimList{}
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return clusterClaims, err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return clusterClaims, err
	}

	l, err := dynamicClient.Resource(helpers.GvrCC).Namespace(cph.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return clusterClaims, err
	}
	for _, ccu := range l.Items {
		cc := &hivev1.ClusterClaim{}
		if runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc); err != nil {
			return clusterClaims, err
		}
		clusterClaims.Items = append(clusterClaims.Items, *cc)
	}
	return clusterClaims, err
}

func getClusterClaimPendingStatus(cc *hivev1.ClusterClaim) *hivev1.ClusterClaimCondition {
	for _, c := range cc.Status.Conditions {
		if c.Type == hivev1.ClusterClaimPendingCondition {
			return &c
		}
	}
	return nil
}

func (cph *ClusterPoolHost) ConvertToPrintClusterClaimList(ccl *hivev1.ClusterClaimList) *printclusterpoolv1alpha1.PrintClusterClaimList {
	pccs := &printclusterpoolv1alpha1.PrintClusterClaimList{}
	for i := range ccl.Items {
		pcc := printclusterpoolv1alpha1.PrintClusterClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ccl.Items[i].Spec.Namespace,
				Namespace: ccl.Items[i].Spec.Namespace,
			},
			Spec: printclusterpoolv1alpha1.PrintClusterClaimSpec{
				ClusterPoolHostName: cph.Name,
				ClusterClaim:        &ccl.Items[i],
				Age:                 helpers.TimeDiff(ccl.Items[i].CreationTimestamp.Time, time.Second),
			},
		}
		clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
		if err != nil {
			pcc.Spec.ErrorMessage = err.Error()
			pccs.Items = append(pccs.Items, pcc)
			continue
		}
		dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
		if err != nil {
			pcc.Spec.ErrorMessage = err.Error()
			pccs.Items = append(pccs.Items, pcc)
			continue
		}
		cdu, err := dynamicClient.Resource(helpers.GvrCD).Namespace(pcc.Spec.ClusterClaim.Spec.Namespace).Get(context.TODO(), pcc.Spec.ClusterClaim.Spec.Namespace, metav1.GetOptions{})
		if err != nil {
			pcc.Spec.ErrorMessage = err.Error()
			pccs.Items = append(pccs.Items, pcc)
			continue
		}
		cd := &hivev1.ClusterDeployment{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd); err != nil {
			pcc.Spec.ErrorMessage = err.Error()
			pccs.Items = append(pccs.Items, pcc)
			continue
		}
		if cd != nil {
			pcc.Spec.PowerState = string(cd.Spec.PowerState)
			pcc.Spec.Hibernate = cd.Labels["hibernate"]
			pcc.Spec.ID = cd.Name
			if ccl.Items[i].Spec.Lifetime != nil {
				pcc.Spec.Lifetime = ccl.Items[i].Spec.Lifetime.Duration.String()
			}
		}
		c := getClusterClaimPendingStatus(pcc.Spec.ClusterClaim)
		if c != nil && c.Status == corev1.ConditionStatus(metav1.ConditionTrue) {
			pcc.Spec.PowerState = string(hivev1.ClusterClaimPendingCondition)
		}
		pathOptions := clientcmd.NewDefaultPathOptions()
		configapi, err := pathOptions.GetStartingConfig()
		if err != nil {
			pcc.Spec.ErrorMessage = err.Error()
			pccs.Items = append(pccs.Items, pcc)
			continue
		}
		pcc.Spec.InUse = configapi.CurrentContext == cph.GetClusterContextName(ccl.Items[i].Name)
		pccs.Items = append(pccs.Items, pcc)
	}
	return pccs
}

func (cph *ClusterPoolHost) GetClusterClaim(clusterName string, waitRunning bool, timeout int, dryRun bool, printFlags *get.PrintFlags) (*hivev1.ClusterClaim, error) {
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return nil, err
	}

	if waitRunning {
		klog.V(3).Infof("Wait cc %s ready", clusterName)
		if err = waitClusterClaimsRunning(dynamicClient, clusterName, "", cph.Namespace, timeout, printFlags); err != nil {
			return nil, err
		}
	}

	ccu, err := dynamicClient.Resource(helpers.GvrCC).Namespace(cph.Namespace).Get(context.TODO(), clusterName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cc := &hivev1.ClusterClaim{}
	if runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc); err != nil {
		return nil, err
	}
	return cc, nil
}

func (cph *ClusterPoolHost) GetPrintClusterClaimCredential(cc *hivev1.ClusterClaim, withCredentials bool) (*printclusterpoolv1alpha1.PrintClusterClaimCredential, error) {
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return nil, err
	}
	if len(cc.Spec.Namespace) == 0 {
		return nil, fmt.Errorf("something wrong happened, the clusterclaim %s doesn't have a spec.namespace set", cc.Name)
	}
	cdu, err := dynamicClient.Resource(helpers.GvrCD).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cd := &hivev1.ClusterDeployment{}
	if runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd); err != nil {
		return nil, err
	}
	ccc := &printclusterpoolv1alpha1.PrintClusterClaimCredential{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cc.Name,
			Namespace: cc.Namespace,
		},
		Spec: printclusterpoolv1alpha1.PrintClusterClaimCredentialSpec{
			User:       "REDACTED",
			Password:   "REDACTED",
			Basedomain: cd.Spec.BaseDomain,
			ApiUrl:     cd.Status.APIURL,
			ConsoleUrl: cd.Status.WebConsoleURL,
		},
	}
	if withCredentials {
		kubeClient, err := kubernetes.NewForConfig(clusterPoolRestConfig)
		if err != nil {
			return nil, err
		}
		s, err := kubeClient.CoreV1().Secrets(cd.Namespace).Get(context.TODO(), cd.Spec.ClusterMetadata.AdminPasswordSecretRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		ccc.Spec.User = string(s.Data["username"])
		ccc.Spec.Password = string(s.Data["password"])
	}
	return ccc, nil
}

func (cph *ClusterPoolHost) PrintClusterClaimCred(cc *hivev1.ClusterClaim, printFlags *get.PrintFlags, withCredentials bool) error {
	cred, err := cph.GetPrintClusterClaimCredential(cc, withCredentials)
	if err != nil {
		return err
	}
	cred.GetObjectKind().
		SetGroupVersionKind(
			schema.GroupVersionKind{
				Group:   printclusterpoolv1alpha1.GroupName,
				Kind:    "PrintClusterClaimCredential",
				Version: printclusterpoolv1alpha1.GroupVersion.Version})
	return helpers.Print(cred, printFlags)
}

func (cph *ClusterPoolHost) OpenClusterClaim(clusterName string, timeout int, printFlags *get.PrintFlags) error {
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}
	klog.V(3).Infof("Wait cc %s ready", clusterName)
	if err = waitClusterClaimsRunning(dynamicClient, clusterName, "", cph.Namespace, timeout, printFlags); err != nil {
		return err
	}

	ccu, err := dynamicClient.Resource(helpers.GvrCC).Namespace(cph.Namespace).Get(context.TODO(), clusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cc := &hivev1.ClusterClaim{}
	if runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc); err != nil {
		return err
	}
	if len(cc.Spec.Namespace) == 0 {
		return fmt.Errorf("something wrong happened, the clusterclaim %s doesn't have a spec.namespace set", cc.Name)
	}
	cdu, err := dynamicClient.Resource(helpers.GvrCD).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cd := &hivev1.ClusterDeployment{}
	if runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd); err != nil {
		return err
	}
	if cd.Spec.PowerState == hivev1.ClusterPowerStateHibernating {
		return fmt.Errorf("%s is hibernating, run a use command to resume it", cc.GetName())
	}
	return helpers.Openbrowser(cd.Status.WebConsoleURL)
}
