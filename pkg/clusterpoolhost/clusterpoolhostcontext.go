// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/stolostron/cm-cli/pkg/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"
)

func (cph *ClusterPoolHost) VerifyClusterPoolContext(
	dryRun bool,
	outputFile string) error {
	token, serviceAccountName, isGlobal, needLogin, err := cph.getClusterPoolSAToken(dryRun, outputFile)
	if err != nil {
		if needLogin {
			fmt.Println(err)
			if errOB := cph.openBrowser(); errOB != nil {
				return errOB
			}
		}
		return err
	}
	return cph.CreateClusterPoolContext(token, serviceAccountName, isGlobal)
}

func (cph *ClusterPoolHost) getClusterPoolSAToken(
	dryRun bool,
	outputFile string) (token, serviceAccountName string, isGlobal, needLogin bool, err error) {
	var clusterPoolRestConfig *rest.Config
	isGlobal = true
	clusterPoolRestConfig, err = cph.GetGlobalRestConfig()
	if err != nil {
		isGlobal = false
		clusterPoolRestConfig, err = GetCurrentRestConfig()
		if err != nil {
			if clusterPoolRestConfig == nil {
				needLogin = true
				err = fmt.Errorf("please login on %s", cph.APIServer)
			}
			return
		}
		if clusterPoolRestConfig.Host != cph.APIServer {
			needLogin = true
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
			needLogin = true
			err = fmt.Errorf("please login on %s", cph.APIServer)
			return
		}

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

func (cph *ClusterPoolHost) openBrowser() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Open browser to %s (Y/N) (default Y): ", cph.Console)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSuffix(answer, "\n")
	klog.V(5).Infof("\nanswer:(%s)\n", answer)
	if strings.ToUpper(answer) == "Y" || len(answer) == 0 {
		return helpers.Openbrowser(cph.Console)
	}
	return nil
}

func (cph *ClusterPoolHost) CreateClusterPoolContext(token, serviceAccountName string, inGlobal bool) error {
	var err error
	var currentConfg *clientcmdapi.Config
	//Get current context
	if inGlobal {
		currentConfg, _, err = GetGlobalConfigAPI()
	} else {
		currentConfg, _, err = GetConfigAPI()

	}
	if err != nil {
		return err
	}
	//Move ClusterPool context
	err = MoveContextToDefault(currentConfg.CurrentContext, cph.GetContextName(), cph.Namespace, serviceAccountName, token)
	return err
}
