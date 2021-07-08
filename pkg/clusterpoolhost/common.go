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

func (c *ClusterPoolHost) VerifyContext(
	dryRun bool,
	outputFile string) error {
	return c.CreateContext(c.Name, dryRun, outputFile, true)
}

func VerifyContext(
	clusterName string,
	dryRun bool,
	outputFile string) error {

	cph, err := GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}
	return cph.CreateContext(clusterName, dryRun, outputFile, false)
}

func (c *ClusterPoolHost) CreateContext(
	clusterName string,
	dryRun bool,
	outputFile string,
	isClusterPool bool) error {
	if isClusterPool {
		return c.setupClusterPool(dryRun, outputFile)
	}
	return c.setupClusterClaim(clusterName, dryRun, outputFile)
}

func (c *ClusterPoolHost) setupClusterPool(
	dryRun bool,
	outputFile string) error {
	inGlobal, err := FindConfigAPIByAPIServer(c.GetContextName(), c.APIServer)
	if err != nil {
		return fmt.Errorf("please login on %s", c.APIServer)
	}
	var clusterPoolRestConfig *rest.Config
	if inGlobal {
		clusterPoolRestConfig, err = GetGlobalCurrentRestConfig()
	} else {
		clusterPoolRestConfig, err = GetCurrentRestConfig()
	}
	if err != nil {
		return err
	}

	//Update the clusterpoolhostfile
	err = c.AddClusterPoolHost(true)
	if err != nil {
		return err
	}

	me, err := c.WhoAmI(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	serviceAccountName := NormalizeName(me.Name)
	//Check if the service account was already created for that user
	//As if already created the me.Name will have this prefix
	if !strings.HasPrefix(me.Name, "system:serviceaccount:"+c.Namespace) {
		err = c.newCKServiceAccount(clusterPoolRestConfig, serviceAccountName, dryRun, outputFile)
		if err != nil {
			return err
		}
	} else {
		serviceAccountName = strings.TrimPrefix(me.Name, "system:serviceaccount:"+c.Namespace+":")
	}

	kubeClient, err := kubernetes.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	// read the token
	token, err := getTokenFromSA(kubeClient, serviceAccountName, c.Namespace)
	if err != nil {
		return err
	}

	return c.CreateClusterPoolContext(token, serviceAccountName, inGlobal)

}

func (c *ClusterPoolHost) setupClusterClaim(
	clusterName string,
	dryRun bool,
	outputFile string) error {
	if err := SetGlobalCurrentContext(c.GetContextName()); err != nil {
		return err
	}

	clusterPoolRestConfig, err := GetGlobalCurrentRestConfig()
	if err != nil {
		return err
	}
	ccRestConfig, err := c.getClusterClaimRestConfig(clusterName, clusterPoolRestConfig)
	if err != nil {
		return err
	}
	kubeClient, err := kubernetes.NewForConfig(ccRestConfig)
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(ccRestConfig)
	if err != nil {
		return err
	}

	apiExtensionsClient, err := apiextensionsclient.NewForConfig(ccRestConfig)
	if err != nil {
		return err
	}

	me, err := c.WhoAmI(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	serviceAccountName := strings.TrimPrefix(me.Name, "system:serviceaccount:"+c.Namespace+":")

	reader := scenario.GetScenarioResourcesReader()

	values := make(map[string]string)
	values["ServiceAccountName"] = serviceAccountName
	output := make([]string, 0)
	files := []string{
		"create/cluster/sa.yaml",
		"create/cluster/cluster-role-binding.yaml",
	}

	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient)
	out, err := applier.ApplyDirectly(reader, values, dryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	if !dryRun {
		token, err := getTokenFromSA(kubeClient, serviceAccountName, "default")
		if err != nil {
			return err
		}
		ccConfigAPI, err := c.getClusterClaimConfigAPI(clusterName, clusterPoolRestConfig)
		if err != nil {
			return err
		}

		return CreateClusterClaimContext(ccConfigAPI, token, clusterName, serviceAccountName)
	} else {
		return clusteradmapply.WriteOutput(outputFile, output)
	}
}

func (c *ClusterPoolHost) newCKServiceAccount(clusterPoolRestConfig *rest.Config, user string, dryRun bool, outputFile string) error {
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
	values["Namespace"] = c.Namespace
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

	if dryRun {
		return clusteradmapply.WriteOutput(outputFile, output)
	}

	//if service-account wait for the sa secret
	err = wait.PollImmediate(1*time.Second, 10*time.Second, func() (bool, error) {
		return waitForSAToken(kubeClient, user, c.Namespace)
	})
	if err != nil {
		return err
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

func (c *ClusterPoolHost) CreateClusterPoolContext(token, serviceAccountName string, inGlobal bool) error {
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
	return MoveContextToDefault(currentContext.CurrentContext, c.GetContextName(), c.Namespace, serviceAccountName, token)
	// if err != nil {
	// 	return err
	// }
	// return SetCurrentContext(c.GetContextName())
}

func CreateClusterClaimContext(configAPI *clientcmdapi.Config, token, clusterName, user string) error {
	return CreateContextFronConfigAPI(configAPI, token, clusterName, DefaultNamespace, user)
}

func (c *ClusterPoolHost) getClusterClaimConfigAPI(clusterName string, clusterPoolRestConfig *rest.Config) (*clientcmdapi.Config, error) {
	kubeClient, err := kubernetes.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return nil, err
	}
	gvr := schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterclaims"}
	ccu, err := dynamicClient.Resource(gvr).Namespace(c.Namespace).Get(context.TODO(), clusterName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cc := &hivev1.ClusterClaim{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc)
	if err != nil {
		return nil, err
	}
	gvr = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterdeployments"}
	cdu, err := dynamicClient.Resource(gvr).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
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

func (c *ClusterPoolHost) getClusterClaimRestConfig(clusterName string, clusterPoolRestConfig *rest.Config) (*rest.Config, error) {
	configapi, err := c.getClusterClaimConfigAPI(clusterName, clusterPoolRestConfig)
	if err != nil {
		return nil, err
	}
	clientConfig := clientcmd.NewDefaultClientConfig(*configapi, nil)
	return clientConfig.ClientConfig()
}
