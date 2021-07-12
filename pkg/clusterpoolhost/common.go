// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost/scenario"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	DefaultNamespace             string = "default"
	ClusterPoolHostContextPrefix string = "clusterpoolhost"
)

var (
	gvrCC = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterclaims"}
	gvrCD = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterdeployments"}
)

func (cph *ClusterPoolHost) VerifyClusterPoolContext(
	dryRun bool,
	outputFile string) error {
	token, serviceAccountName, isGlobal, err := cph.getClusterPoolSAToken(dryRun, outputFile)
	if err != nil {
		return err
	}
	return cph.CreateClusterPoolContext(token, serviceAccountName, isGlobal)
}

func VerifyClusterClaimContext(
	clusterName string,
	timeout int,
	dryRun bool,
	outputFile string) error {

	cph, err := GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}
	token, serviceAccountName, ccConfigAPI, err := cph.getClusterClaimSAToken(clusterName, timeout, dryRun, outputFile)
	if err != nil {
		return err
	}

	return CreateClusterClaimContext(ccConfigAPI, token, clusterName, serviceAccountName)
}

func (cph *ClusterPoolHost) getClusterPoolSAToken(
	dryRun bool,
	outputFile string) (token, serviceAccountName string, isGlobal bool, err error) {
	isGlobal = true
	err = SetCPHContext(cph.GetContextName())
	if err != nil {
		isGlobal, err = findConfigAPIByAPIServer(cph.GetContextName(), cph.APIServer)
		if err != nil {
			err = fmt.Errorf("please login on %s", cph.APIServer)
			return
		}
	}
	var clusterPoolRestConfig *rest.Config
	if isGlobal {
		clusterPoolRestConfig, err = GetGlobalCurrentRestConfig()
	} else {
		clusterPoolRestConfig, err = GetCurrentRestConfig()
	}
	if err != nil {
		return
	}

	//Update the clusterpoolhostfile
	err = cph.AddClusterPoolHost()
	if err != nil {
		return
	}

	me, err := WhoAmI(clusterPoolRestConfig)
	if err != nil {
		return
	}

	serviceAccountName = NormalizeName(me.Name)
	//Check if the service account was already created for that user
	//As if already created the me.Name will have this prefix
	if !strings.HasPrefix(me.Name, "system:serviceaccount:"+cph.Namespace) {
		err = cph.newCKServiceAccount(clusterPoolRestConfig, serviceAccountName, dryRun, outputFile)
		if err != nil {
			return
		}
	} else {
		serviceAccountName = strings.TrimPrefix(me.Name, "system:serviceaccount:"+cph.Namespace+":")
	}

	kubeClient, err := kubernetes.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return
	}

	// read the token
	token, err = getTokenFromSA(kubeClient, serviceAccountName, cph.Namespace)
	return
}

func (cph *ClusterPoolHost) getClusterClaimSAToken(
	clusterName string,
	timeout int,
	dryRun bool,
	outputFile string) (token, serviceAccountName string, ccConfigAPI *clientcmdapi.Config, err error) {
	if err = SetGlobalCurrentContext(cph.GetContextName()); err != nil {
		return
	}

	clusterPoolRestConfig, err := GetGlobalCurrentRestConfig()
	if err != nil {
		return
	}

	dynamicClientCP, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return
	}

	me, err := WhoAmI(clusterPoolRestConfig)
	if err != nil {
		return
	}

	serviceAccountName = strings.TrimPrefix(me.Name, "system:serviceaccount:"+cph.Namespace+":")

	reader := scenario.GetScenarioResourcesReader()

	values := make(map[string]string)
	values["ServiceAccountName"] = serviceAccountName
	output := make([]string, 0)

	files := []string{
		"create/cluster/sa.yaml",
		"create/cluster/cluster-role-binding.yaml",
	}

	applierBuilder := &clusteradmapply.ApplierBuilder{}
	if !dryRun {
		if err = setHibernateClusterClaims(clusterName, false, "", dryRun, outputFile); err != nil {
			return
		}
		if err = waitClusterClaimsRunning(dynamicClientCP, clusterName, "", cph.Namespace, timeout); err != nil {
			return
		}
		ccRestConfig, errG := cph.getClusterClaimRestConfig(clusterName, clusterPoolRestConfig)
		if errG != nil {
			err = errG
			return
		}
		kubeClientCC, errG := kubernetes.NewForConfig(ccRestConfig)
		if err != nil {
			err = errG
			return
		}

		dynamicClientCC, errG := dynamic.NewForConfig(ccRestConfig)
		if err != nil {
			err = errG
			return
		}

		apiExtensionsClientCC, errG := apiextensionsclient.NewForConfig(ccRestConfig)
		if err != nil {
			err = errG
			return
		}

		applier := applierBuilder.WithClient(kubeClientCC, apiExtensionsClientCC, dynamicClientCC)
		out, errG := applier.ApplyDirectly(reader, values, dryRun, "", files...)
		if err != nil {
			err = errG
			return
		}
		output = append(output, out...)
		token, err = getTokenFromSA(kubeClientCC, serviceAccountName, "default")
		if err != nil {
			return
		}
		ccConfigAPI, err = cph.getClusterClaimConfigAPI(clusterName, clusterPoolRestConfig)
		if err != nil {
			return
		}
	} else {
		applier := applierBuilder
		out, errG := applier.MustTemplateAssets(reader, values, "", files...)
		if err != nil {
			err = errG
			return
		}
		output = append(output, out...)
	}

	err = clusteradmapply.WriteOutput(outputFile, output)
	if err != nil {
		return
	}

	if !dryRun {
	}

	return
}

func (cph *ClusterPoolHost) newCKServiceAccount(clusterPoolRestConfig *rest.Config, user string, dryRun bool, outputFile string) error {
	reader := scenario.GetScenarioResourcesReader()

	kubeClient, err := kubernetes.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	apiExtensionsClient, err := apiextensionsclient.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	values := make(map[string]string)
	values["Name"] = user
	values["Namespace"] = cph.Namespace
	output := make([]string, 0)
	files := []string{
		"create/clusterpoolhost/sa.yaml",
	}
	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient)
	out, err := applier.ApplyDirectly(reader, values, dryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	err = clusteradmapply.WriteOutput(outputFile, output)
	if err != nil {
		return err
	}

	if !dryRun {
		//if service-account wait for the sa secret
		err = wait.PollImmediate(1*time.Second, 10*time.Second, func() (bool, error) {
			return waitForSAToken(kubeClient, user, cph.Namespace)
		})
		if err != nil {
			return err
		}
	}
	return nil

}

func waitForSAToken(kubeClient kubernetes.Interface, serviceAccountName, namespace string) (bool, error) {
	_, err := getTokenFromSA(kubeClient, serviceAccountName, namespace)
	switch {
	case errors.IsNotFound(err):
		return false, nil
	case err != nil:
		return false, err
	}
	return true, nil
}

func getSecretFromSA(
	kubeClient kubernetes.Interface,
	serviceAccountName, namespace string) (*corev1.Secret, error) {
	sa, err := kubeClient.CoreV1().
		ServiceAccounts(namespace).
		Get(context.TODO(), serviceAccountName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	var secret *corev1.Secret
	var prefix string
	for _, objectRef := range sa.Secrets {
		if objectRef.Namespace != "" && objectRef.Namespace != namespace {
			continue
		}
		prefix = serviceAccountName
		if len(prefix) > 63 {
			prefix = prefix[:37]
		}
		if strings.HasPrefix(objectRef.Name, prefix) {
			secret, err = kubeClient.CoreV1().
				Secrets(namespace).
				Get(context.TODO(), objectRef.Name, metav1.GetOptions{})
			if err != nil {
				continue
			}
			if secret.Type == corev1.SecretTypeServiceAccountToken {
				break
			}
		}
	}
	if secret == nil {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    corev1.GroupName,
			Resource: "secrets"},
			fmt.Sprintf("secret with prefix %s and type %s not found in service account %s/%s",
				prefix,
				corev1.SecretTypeServiceAccountToken,
				namespace,
				serviceAccountName))
	}
	return secret, nil
}

//GetBootstrapSecretFromSA retrieves the service-account token secret
func getTokenFromSA(kubeClient kubernetes.Interface, serviceAccountName, namespace string) (string, error) {
	secret, err := getSecretFromSA(kubeClient, serviceAccountName, namespace)
	if err != nil {
		return "", err
	}
	return string(secret.Data["token"]), nil
}

func (cph *ClusterPoolHost) CreateClusterPoolContext(token, serviceAccountName string, inGlobal bool) error {
	var err error
	var currentContext *clientcmdapi.Config
	//Get current context
	if inGlobal {
		currentContext, err = GetGlobalConfig()
	} else {
		currentContext, err = GetConfig()

	}
	if err != nil {
		return err
	}

	//Move ClusterPool context
	return MoveContextToDefault(currentContext.CurrentContext, cph.GetContextName(), cph.Namespace, serviceAccountName, token)
}

func CreateClusterClaimContext(configAPI *clientcmdapi.Config, token, clusterName, user string) error {
	return CreateContextFronConfigAPI(configAPI, token, clusterName, DefaultNamespace, user)
}

func (cph *ClusterPoolHost) getClusterClaimConfigAPI(clusterName string, clusterPoolRestConfig *rest.Config) (*clientcmdapi.Config, error) {
	kubeClient, err := kubernetes.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return nil, err
	}
	ccu, err := dynamicClient.Resource(gvrCC).Namespace(cph.Namespace).Get(context.TODO(), clusterName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cc := &hivev1.ClusterClaim{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc)
	if err != nil {
		return nil, err
	}
	cdu, err := dynamicClient.Resource(gvrCD).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cd := &hivev1.ClusterDeployment{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd)
	if err != nil {
		return nil, err
	}
	s, err := kubeClient.CoreV1().Secrets(cd.Namespace).Get(context.TODO(), cd.Spec.ClusterMetadata.AdminKubeconfigSecretRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return clientcmd.Load(s.Data["kubeconfig"])
}

func (cph *ClusterPoolHost) getClusterClaimRestConfig(clusterName string, clusterPoolRestConfig *rest.Config) (*rest.Config, error) {
	configapi, err := cph.getClusterClaimConfigAPI(clusterName, clusterPoolRestConfig)
	if err != nil {
		return nil, err
	}
	clientConfig := clientcmd.NewDefaultClientConfig(*configapi, nil)
	return clientConfig.ClientConfig()
}

func CreateClusterClaims(clusterClaimNames, clusterPoolName string, skipSchedule bool, timeout int, dryRun bool, outputFile string) error {
	cph, err := GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}
	err = cph.VerifyClusterPoolContext(dryRun, outputFile)
	if err != nil {
		return err
	}
	clusterPoolRestConfig, err := GetGlobalCurrentRestConfig()
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

	gvr := schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterpools"}
	_, err = dynamicClient.Resource(gvr).Namespace(cph.Namespace).Get(context.TODO(), clusterPoolName, metav1.GetOptions{})
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
	if err = RunClusterClaims(clusterClaimNames, skipSchedule, dryRun, outputFile); err != nil {
		return err
	}
	return clusteradmapply.WriteOutput(outputFile, output)
}

func RunClusterClaims(clusterClaimNames string, skipSchedule, dryRun bool, outputFile string) error {
	skipScheduleAction := "true"
	if skipSchedule {
		skipScheduleAction = "skip"
	}
	return setHibernateClusterClaims(clusterClaimNames, false, skipScheduleAction, dryRun, outputFile)
}

func HibernateClusterClaims(clusterClaimNames string, skipSchedule, dryRun bool, outputFile string) error {
	skipScheduleAction := "true"
	if skipSchedule {
		skipScheduleAction = "skip"
	}
	return setHibernateClusterClaims(clusterClaimNames, true, skipScheduleAction, dryRun, outputFile)
}

func setHibernateClusterClaims(clusterClaimNames string, hibernate bool, skipScheduleAction string, dryRun bool, outputFile string) error {
	cph, err := GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}
	err = cph.VerifyClusterPoolContext(dryRun, outputFile)
	if err != nil {
		return err
	}
	clusterPoolRestConfig, err := GetGlobalCurrentRestConfig()
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
		gvr := schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterpools"}
		cpu, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), clusterPoolName, metav1.GetOptions{})
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
		for _, c := range cc.Status.Conditions {
			if c.Type == hivev1.ClusterRunningCondition &&
				c.Status == corev1.ConditionStatus(metav1.ConditionTrue) {
				running = true
				fmt.Printf("clusterclaim %s is running with id %s (%d/%d)\n", clusterClaimName, cc.Spec.Namespace, i, timeout)
				break
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

func DeleteClusterClaims(clusterClaimNames string, dryRun bool, outputFile string) error {
	cph, err := GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}
	err = cph.VerifyClusterPoolContext(dryRun, outputFile)
	if err != nil {
		return err
	}
	clusterPoolRestConfig, err := GetGlobalCurrentRestConfig()
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
	err = cph.VerifyClusterPoolContext(dryRun, outputFile)
	if err != nil {
		return err
	}
	clusterPoolRestConfig, err := GetGlobalCurrentRestConfig()
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
	if len(l.Items) == 0 {
		fmt.Printf("No clusterclaim found for clusterpoolhost %s\n", cph.Name)
	}
	for _, ccu := range l.Items {
		cc := &hivev1.ClusterClaim{}
		if runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc); err != nil {
			return err
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
		if showCphName {
			fmt.Printf("%-15s\t%-15s\t%-11s\t%-4s\twith id %s\n", cph.Name, cc.GetName(), cd.Spec.PowerState, cd.Labels["hibernate"], cd.GetName())
		} else {
			fmt.Printf("%-15s\t%-11s\t%-4s\twith id %s\n", cc.GetName(), cd.Spec.PowerState, cd.Labels["hibernate"], cd.GetName())
		}
	}
	return nil
}

func GetClusterClaim(clusterName string, dryRun bool, outputFile string) error {
	cph, err := GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}
	err = cph.VerifyClusterPoolContext(dryRun, outputFile)
	if err != nil {
		return err
	}
	clusterPoolRestConfig, err := GetGlobalCurrentRestConfig()
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
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
