// Copyright Contributors to the Open Cluster Management project
package contexts

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"
	"github.com/stolostron/cm-cli/pkg/managedcluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clusterv1 "open-cluster-management.io/api/cluster/v1"

	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	return nil
}

func (o *Options) validate() error {
	return nil
}

func (o *Options) run(streams genericclioptions.IOStreams) (err error) {

	cphs, err := clusterpoolhost.GetClusterPoolHosts()
	if err != nil {
		return err
	}
	var cph *clusterpoolhost.ClusterPoolHost
	if len(cphs.ClusterPoolHosts) != 0 {
		cph, err = clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
		if err != nil {
			fmt.Println("no clusterpoolhost found, will only get the contexts of hive generated clusters")
		}
	}
	restConfig, err := o.CMFlags.KubectlFactory.ToRESTConfig()
	if err != nil {
		return err
	}

	clusterClient, err := clusterclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	mcs, err := clusterClient.ClusterV1().ManagedClusters().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	cmdAPIConfig := &clientcmdapi.Config{
		Kind:       "Config",
		APIVersion: "v1",
		AuthInfos:  make(map[string]*clientcmdapi.AuthInfo),
		Contexts:   make(map[string]*clientcmdapi.Context),
		Clusters:   make(map[string]*clientcmdapi.Cluster),
	}

	currentCmdAPIConfig, _, err := clusterpoolhost.GetConfigAPI()
	if err != nil {
		return err
	}
	cmdAPIConfig.AuthInfos[currentCmdAPIConfig.CurrentContext] = currentCmdAPIConfig.AuthInfos[currentCmdAPIConfig.CurrentContext]
	cmdAPIConfig.Contexts[currentCmdAPIConfig.CurrentContext] = currentCmdAPIConfig.Contexts[currentCmdAPIConfig.CurrentContext]
	cmdAPIConfig.Clusters[currentCmdAPIConfig.CurrentContext] = currentCmdAPIConfig.Clusters[currentCmdAPIConfig.CurrentContext]
	cmdAPIConfig.CurrentContext = currentCmdAPIConfig.CurrentContext

	dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}

	kubeClient, err := o.CMFlags.KubectlFactory.KubernetesClientSet()
	if err != nil {
		return err
	}

	for _, mc := range mcs.Items {
		if mc.Name == "local-cluster" {
			continue
		}
		clusterCmdAPIConfig, err := managedcluster.GetCmdAPIConfig(dynamicClient, kubeClient, &mc, cph)
		if err != nil {
			return err
		}
		if clusterCmdAPIConfig == nil {
			fmt.Fprintf(streams.ErrOut, "no kubeconfig found for managedcluster %s\n", mc.Name)
		}
		addCluster(&mc, cmdAPIConfig, clusterCmdAPIConfig)
	}
	// return err
	data, err := clientcmd.Write(*cmdAPIConfig)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func addCluster(mc *clusterv1.ManagedCluster, configs, config *clientcmdapi.Config) {
	if config != nil {
		for _, v := range config.AuthInfos {
			configs.AuthInfos[mc.Name] = v
		}
		for _, v := range config.Clusters {
			configs.Clusters[mc.Name] = v
		}
		for _, v := range config.Contexts {
			v.AuthInfo = mc.Name
			v.Cluster = mc.Name
			configs.Contexts[mc.Name] = v
		}
	}
}
