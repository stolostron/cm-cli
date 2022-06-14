// Copyright Contributors to the Open Cluster Management project
package hypershiftdeployment

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/cmd/get/config/hypershiftdeployment/scenario"
	"github.com/stolostron/cm-cli/pkg/hypershift"
	apiextensionsClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 0 {
		o.clusterName = args[0]
	}
	return nil
}

func (o *Options) validate() error {
	if len(o.clusterName) == 0 {
		return fmt.Errorf("hypershiftdeployment name is missing")
	}
	return nil
}

func (o *Options) run() (err error) {
	kubeClient, apiextensionsClient, dynamicClient, err := clusteradmhelpers.GetClients(o.CMFlags.KubectlFactory)
	if err != nil {
		return err
	}
	return o.runWithClient(kubeClient, apiextensionsClient, dynamicClient)
}

func (o *Options) runWithClient(kubeClient kubernetes.Interface,
	apiExtensionsClient apiextensionsClient.Interface,
	dynamicClient dynamic.Interface) (err error) {
	reader := scenario.GetScenarioResourcesReader()

	//Get hypershiftdeployment
	hd, err := hypershift.GetHypershiftDeployment(o.clusterName, o.CMFlags)
	if err != nil {
		return err
	}

	klog.V(5).Infof("%v\n", hd)
	applierBuilder := clusteradmapply.NewApplierBuilder()
	applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient).Build()
	b, err := applier.MustTemplateAsset(reader, *hd, "", "config/config.yaml")
	if err != nil {
		return err
	}
	if len(o.outputFile) != 0 {
		return ioutil.WriteFile(o.outputFile, b, 0600)
	}
	fmt.Printf("%s\n", string(b))
	return nil
}
