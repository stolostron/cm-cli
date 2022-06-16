// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"context"
	"fmt"
	"strings"
	"time"

	userv1 "github.com/openshift/api/user/v1"
	userv1typedclient "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	"github.com/stolostron/cm-cli/pkg/clusterpoolhost/scenario"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"
)

//WhoAmI returns the current user
func WhoAmI(restConfig *rest.Config) (*userv1.User, error) {
	userInterface, err := userv1typedclient.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	me, err := userInterface.Users().Get(context.TODO(), "~", metav1.GetOptions{})
	if err == nil {
		return me, err
	}

	return me, err
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
		"create/clusterpoolhost/secret-token.yaml",
	}
	applierBuilder := clusteradmapply.NewApplierBuilder()
	applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient).Build()
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
