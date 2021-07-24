// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"context"
	"fmt"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/attach/cluster/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	//Check if default values must be used
	if o.valuesPath == "" {
		if len(args) < 1 {
			return fmt.Errorf("ClusterClaim name is missing")
		}
		o.ClusterClaim = args[0]
		reader := scenario.GetScenarioResourcesReader()
		o.values, err = helpers.ConvertReaderFileToValuesMap(valuesDefaultPath, reader)
		if err != nil {
			return err
		}
		mc := o.values["managedCluster"].(map[string]interface{})
		mc["name"] = o.ClusterClaim
	} else {
		//Read values
		o.values, err = helpers.ConvertValuesFileToValuesMap(o.valuesPath, "")
		if err != nil {
			return err
		}
	}

	imc, ok := o.values["managedCluster"]
	if !ok || imc == nil {
		return fmt.Errorf("managedCluster is missing")
	}
	mc := imc.(map[string]interface{})

	if _, ok := mc["labels"]; !ok {
		mc["labels"] = map[string]interface{}{
			"cloud":  "auto-detect",
			"vendor": "auto-detect",
		}
	}

	ilabels := mc["labels"]
	labels := ilabels.(map[string]interface{})
	if _, ok := labels["vendor"]; !ok {
		labels["vendor"] = "auto-detect"
	}

	if _, ok := labels["cloud"]; !ok {
		labels["cloud"] = "auto-detect"
	}

	return nil
}

func (o *Options) validate() error {
	imc, ok := o.values["managedCluster"]
	if !ok || imc == nil {
		return fmt.Errorf("managedCluster is missing")
	}
	mc := imc.(map[string]interface{})

	if o.ClusterClaim == "" {
		iname, ok := mc["name"]
		if !ok || iname == nil {
			return fmt.Errorf("cluster name is missing")
		}
		o.ClusterClaim = iname.(string)
		if len(o.ClusterClaim) == 0 {
			return fmt.Errorf("managedCluster.name not specified")
		}
	}

	mc["name"] = o.ClusterClaim

	return nil
}

func (o *Options) run() (err error) {
	cphs, err := clusterpoolhost.GetClusterPoolHosts()
	if err != nil {
		return err
	}

	currentCph, err := cphs.GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}

	if len(o.ClusterPoolHost) != 0 {
		cph, err := cphs.GetClusterPoolHost(o.ClusterPoolHost)
		if err != nil {
			return err
		}

		err = cphs.SetActive(cph)
		if err != nil {
			return err
		}
	}

	err = o.attachClusterClaim(cphs)

	if len(o.ClusterPoolHost) != 0 {
		if err := cphs.SetActive(currentCph); err != nil {
			return err
		}
	}
	return err

}

func (o *Options) attachClusterClaim(cphs *clusterpoolhost.ClusterPoolHosts) error {

	output := make([]string, 0)
	reader := scenario.GetScenarioResourcesReader()

	files := []string{
		"attach/hub/namespace.yaml",
		"attach/hub/managed_cluster_secret.yaml",
	}

	cph, err := clusterpoolhost.GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}

	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return err
	}

	dynamicClientCP, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	ccu, err := dynamicClientCP.Resource(helpers.GvrCC).Namespace(cph.Namespace).Get(context.TODO(), o.ClusterClaim, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cc := &hivev1.ClusterClaim{}
	if err = runtime.DefaultUnstructuredConverter.FromUnstructured(ccu.UnstructuredContent(), cc); err != nil {
		return err
	}
	cdu, err := dynamicClientCP.Resource(helpers.GvrCD).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cd := &hivev1.ClusterDeployment{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd)
	if err != nil {
		return err
	}

	kubeClientCP, err := kubernetes.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	kubeConfigSecret, err := kubeClientCP.CoreV1().
		Secrets(cd.GetNamespace()).
		Get(context.TODO(), cd.Spec.ClusterMetadata.AdminKubeconfigSecretRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	imc := o.values["managedCluster"]
	mc := imc.(map[string]interface{})
	mc["kubeConfig"] = string(kubeConfigSecret.Data["kubeconfig"])
	klog.V(5).Infof("KubeConfig:\n%s\n", kubeConfigSecret.Data["kubeconfig"])

	hubRestConfig, err := clusterpoolhost.GetCurrentRestConfig()
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(hubRestConfig)
	if err != nil {
		return err
	}

	constraint := ">=2.3.0"
	supported, err := helpers.IsSupported(kubeClient, constraint)
	if err != nil {
		return err
	}
	if !supported {
		return fmt.Errorf("this command requires RHACM version %s", constraint)
	}

	dynamicClient, err := dynamic.NewForConfig(hubRestConfig)
	if err != nil {
		return err
	}

	apiExtensionsClient, err := apiextensionsclient.NewForConfig(hubRestConfig)
	if err != nil {
		return err
	}

	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient).Build()
	out, err := applier.ApplyDirectly(reader, o.values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	files = []string{
		"attach/hub/managed_cluster_cr.yaml",
		"attach/hub/klusterlet_addon_config_cr.yaml",
	}

	out, err = applier.ApplyCustomResources(reader, o.values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	return clusteradmapply.WriteOutput(o.outputFile, output)
}
