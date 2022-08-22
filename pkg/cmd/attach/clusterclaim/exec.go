// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"context"
	"fmt"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
	workclientset "open-cluster-management.io/api/client/work/clientset/versioned"

	"github.com/spf13/cobra"
	"github.com/stolostron/applier/pkg/apply"
	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"
	"github.com/stolostron/cm-cli/pkg/cmd/attach/cluster/scenario"
	"github.com/stolostron/cm-cli/pkg/helpers"
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
	cph, err := clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	err = o.attachClusterClaim(cph)

	return err

}

func (o *Options) attachClusterClaim(cph *clusterpoolhost.ClusterPoolHost) error {

	output := make([]string, 0)
	reader := scenario.GetScenarioResourcesReader()

	files := []string{
		"attach/hub/namespace.yaml",
		"attach/hub/managed_cluster_secret.yaml",
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

	rhacmConstraint := ">=2.3.0"
	mceConstraint := ">=1.0.0"

	supported, platform, err := helpers.IsSupportedVersion(o.CMFlags, true, o.ClusterPoolHost, rhacmConstraint, mceConstraint)
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

	applierBuilder := apply.NewApplierBuilder()
	applier := applierBuilder.WithRestConfig(hubRestConfig).Build()
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
			return helpers.WaitKlusterlet(clusterClient, o.ClusterClaim, o.timeout)
		}
		if o.waitAddOns {
			return helpers.WaitKlusterletAddons(workClient, o.ClusterClaim, o.timeout)
		}
	}

	return apply.WriteOutput(o.outputFile, output)
}
