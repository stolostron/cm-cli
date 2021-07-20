// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
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

func (cph *ClusterPoolHost) getClusterPoolSAToken(
	dryRun bool,
	outputFile string) (token, serviceAccountName string, isGlobal bool, err error) {
	var clusterPoolRestConfig *rest.Config
	isGlobal = true
	err = SetCPHContext(cph.GetContextName())
	if err != nil {
		clusterPoolRestConfig, err = GetCurrentRestConfig()
		if err != nil {
			if clusterPoolRestConfig == nil {
				err = fmt.Errorf("please login on %s", cph.APIServer)
			}
			return
		}
		if clusterPoolRestConfig.Host != cph.APIServer {
			err = fmt.Errorf("please login on %s", cph.APIServer)
			return
		}
		var kubeClient kubernetes.Interface
		kubeClient, err = kubernetes.NewForConfig(clusterPoolRestConfig)
		if err != nil {
			return
		}
		_, err = kubeClient.CoreV1().Secrets(cph.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			err = fmt.Errorf("please login on %s", cph.APIServer)
			return
		}
	}
	clusterPoolRestConfig, err = GetGlobalCurrentRestConfig()
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

func (cph *ClusterPoolHost) CreateClusterPoolContext(token, serviceAccountName string, inGlobal bool) error {
	var err error
	var currentContext *clientcmdapi.Config
	//Get current context
	if inGlobal {
		currentContext, _, err = GetGlobalConfigAPI()
	} else {
		currentContext, _, err = GetConfigAPI()

	}
	if err != nil {
		return err
	}

	//Move ClusterPool context
	return MoveContextToDefault(currentContext.CurrentContext, cph.GetContextName(), cph.Namespace, serviceAccountName, token)
}
