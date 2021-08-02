// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"context"
	"strings"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"
)

const (
	DefaultNamespace string = "default"
)

func UseClusterClaimContext(
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

func (cph *ClusterPoolHost) getClusterClaimSAToken(
	clusterName string,
	timeout int,
	dryRun bool,
	outputFile string) (token, serviceAccountName string, ccConfigAPI *clientcmdapi.Config, err error) {

	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
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
		if err = cph.setHibernateClusterClaims(clusterName, false, "", dryRun, outputFile); err != nil {
			return
		}
		if err = waitClusterClaimsRunning(dynamicClientCP, clusterName, "", cph.Namespace, timeout, nil); err != nil {
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

	return
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
	ccu, err := dynamicClient.Resource(helpers.GvrCC).Namespace(cph.Namespace).Get(context.TODO(), clusterName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cc := &hivev1.ClusterClaim{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc)
	if err != nil {
		return nil, err
	}
	cdu, err := dynamicClient.Resource(helpers.GvrCD).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
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
