// Copyright Contributors to the Open Cluster Management project
package hostedcluster

import (
	"context"
	"fmt"

	hypershiftv1alpha1 "github.com/openshift/hypershift/api/v1alpha1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
	workclientset "open-cluster-management.io/api/client/work/clientset/versioned"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/cmd/attach/cluster/scenario"
	"github.com/stolostron/cm-cli/pkg/helpers"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	//Check if default values must be used
	if o.valuesPath == "" {
		if len(args) < 1 {
			return fmt.Errorf("ClusterClaim name is missing")
		}
		o.HostedCluster = args[0]
		reader := scenario.GetScenarioResourcesReader()
		o.values, err = helpers.ConvertReaderFileToValuesMap(valuesDefaultPath, reader)
		if err != nil {
			return err
		}
		mc := o.values["managedCluster"].(map[string]interface{})
		mc["name"] = o.HostedCluster
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
	dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}
	return o.validateWithClient(dynamicClient)
}

func (o *Options) validateWithClient(dynamicClient dynamic.Interface) error {
	imc, ok := o.values["managedCluster"]
	if !ok || imc == nil {
		return fmt.Errorf("managedCluster is missing")
	}
	mc := imc.(map[string]interface{})

	if o.HostedCluster == "" {
		iname, ok := mc["name"]
		if !ok || iname == nil {
			return fmt.Errorf("cluster name is missing")
		}
		o.HostedCluster = iname.(string)
		if len(o.HostedCluster) == 0 {
			return fmt.Errorf("managedCluster.name not specified")
		}
	}

	mc["name"] = o.HostedCluster

	if _, err := dynamicClient.Resource(helpers.GvrHC).Namespace(o.HostedClusterNamespace).Get(context.TODO(), o.HostedCluster, metav1.GetOptions{}); err != nil {
		return fmt.Errorf("%s is not a hostedcluster, %s", o.HostedCluster, err)
	}
	return nil
}

func (o *Options) run() (err error) {
	return o.attachHostedCluster()
}

func (o *Options) attachHostedCluster() error {

	output := make([]string, 0)
	reader := scenario.GetScenarioResourcesReader()

	files := []string{
		"attach/hub/namespace.yaml",
		"attach/hub/managed_cluster_secret.yaml",
	}

	kubeClient, err := o.CMFlags.KubectlFactory.KubernetesClientSet()
	if err != nil {
		return err
	}
	dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}

	hcu, err := dynamicClient.Resource(helpers.GvrHC).Namespace(o.HostedClusterNamespace).Get(context.TODO(), o.HostedCluster, metav1.GetOptions{})
	if err != nil {
		return err
	}
	hc := &hypershiftv1alpha1.HostedCluster{}
	if err = runtime.DefaultUnstructuredConverter.FromUnstructured(hcu.UnstructuredContent(), hc); err != nil {
		return err
	}
	if hc.Status.KubeConfig == nil {
		return fmt.Errorf("kubeconfig not yet available for cluster %s", o.HostedCluster)
	}
	kubeConfigSecret, err := kubeClient.CoreV1().Secrets(o.HostedClusterNamespace).Get(context.TODO(), hc.Status.KubeConfig.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	imc := o.values["managedCluster"]
	mc := imc.(map[string]interface{})
	mc["kubeConfig"] = string(kubeConfigSecret.Data["kubeconfig"])

	klog.V(5).Infof("KubeConfig:\n%s\n", kubeConfigSecret.Data["kubeconfig"])

	rhacmConstraint := ">=2.3.0"
	mceConstraint := ">=1.0.0"

	supported, platform, err := helpers.IsSupportedVersion(o.CMFlags, true, o.HostedCluster, rhacmConstraint, mceConstraint)
	if err != nil {
		return err
	}
	if !supported {
		switch platform {
		case helpers.RHACM:
			return fmt.Errorf("this command requires %s version %s", platform, rhacmConstraint)
		case helpers.MCE:
			return fmt.Errorf("this command requires %s version %s", platform, mceConstraint)
		}
	}

	hubRestConfig, err := o.CMFlags.KubectlFactory.ToRESTConfig()
	if err != nil {
		return err
	}

	apiExtensionsClient, err := apiextensionsclient.NewForConfig(hubRestConfig)
	if err != nil {
		return err
	}

	applierBuilder := clusteradmapply.NewApplierBuilder()
	applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient).Build()
	out, err := applier.ApplyDirectly(reader, o.values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	files = []string{
		"attach/hub/managed_cluster_cr.yaml",
	}

	if helpers.IsRHACM(o.CMFlags) {
		files = append(files, "attach/hub/klusterlet_addon_config_cr.yaml")
	}

	out, err = applier.ApplyCustomResources(reader, o.values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	if !o.CMFlags.DryRun {
		clusterClient, err := clusterclientset.NewForConfig(hubRestConfig)
		if err != nil {
			return err
		}
		workClient, err := workclientset.NewForConfig(hubRestConfig)
		if err != nil {
			return err
		}
		if o.waitAgent || o.waitAddOns {
			return helpers.WaitKlusterlet(clusterClient, o.HostedCluster, o.timeout)
		}
		if o.waitAddOns {
			return helpers.WaitKlusterletAddons(workClient, o.HostedCluster, o.timeout)
		}
	}

	return clusteradmapply.WriteOutput(o.outputFile, output)
}
