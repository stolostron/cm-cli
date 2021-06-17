// Copyright Contributors to the Open Cluster Management project
package addons

import (
	"fmt"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"

	attachscenario "github.com/open-cluster-management/cm-cli/pkg/cmd/attach/cluster/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/enable/addons/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	//Check if default values must be used
	if o.valuesPath == "" {
		if o.clusterName != "" {
			reader := scenario.GetScenarioResourcesReader()
			o.values, err = helpers.ConvertReaderFileToValuesMap(valuesDefaultPath, reader)
			if err != nil {
				return err
			}
			mc := o.values["managedCluster"].(map[string]interface{})
			mc["name"] = o.clusterName
		} else {
			return fmt.Errorf("values or name are missing")
		}
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

	return nil
}

func (o *Options) validate() error {
	kubeClient, err := o.CMFlags.KubectlFactory.KubernetesClientSet()
	if err != nil {
		return err
	}
	dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}
	return o.validateWithClient(kubeClient, dynamicClient)
}

func (o *Options) validateWithClient(kubeClient kubernetes.Interface, dynamicClient dynamic.Interface) error {
	imc, ok := o.values["managedCluster"]
	if !ok || imc == nil {
		return fmt.Errorf("managedCluster is missing")
	}
	mc := imc.(map[string]interface{})

	if o.clusterName == "" {
		iname, ok := mc["name"]
		if !ok || iname == nil {
			return fmt.Errorf("cluster name is missing")
		}
		o.clusterName = iname.(string)
		if len(o.clusterName) == 0 {
			return fmt.Errorf("managedCluster.name not specified")
		}
	}

	mc["name"] = o.clusterName

	return nil
}

func (o *Options) run() (err error) {
	dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}
	discoveryClient, err := o.CMFlags.KubectlFactory.ToDiscoveryClient()
	if err != nil {
		return err
	}
	return o.runWithClient(dynamicClient, discoveryClient)
}

func (o *Options) runWithClient(dynamicClient dynamic.Interface,
	discoveryClient discovery.DiscoveryInterface) (err error) {
	output := make([]string, 0)
	reader := attachscenario.GetScenarioResourcesReader()

	files := []string{
		"attach/hub/klusterlet_addon_config_cr.yaml",
	}

	out, err := clusteradmapply.ApplyCustomResouces(dynamicClient, discoveryClient, reader, o.values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)
	return clusteradmapply.WriteOutput(o.outputFile, output)
}
