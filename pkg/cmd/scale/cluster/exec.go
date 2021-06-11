// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"

	"github.com/open-cluster-management/cm-cli/pkg/cmd/scale/cluster/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	var mc map[string]interface{}
	if o.valuesPath == "" {
		reader := scenario.GetScenarioResourcesReader()
		o.values, err = helpers.ConvertReaderFileToValuesMap(valuesDefaultPath, reader)
		if err != nil {
			return err
		}
	} else {
		o.values, err = helpers.ConvertValuesFileToValuesMap(o.valuesPath, "")
		if err != nil {
			return err
		}
	}
	if imc, ok := o.values["managedCluster"]; !ok {
		return fmt.Errorf("managedCluster is missing")
	} else {
		mc = imc.(map[string]interface{})
	}
	// overwrite values with parameters
	if o.clusterName != "" {
		mc["name"] = o.clusterName
	}
	if o.machinePoolName != "" {
		mc["machinepool"] = o.machinePoolName
	}
	if cmd.Flags().Changed("replicas") {
		mc["replicas"] = o.replicas
	}

	//Align parameters with values
	if mc["name"] != nil {
		o.clusterName = mc["name"].(string)
	}
	if mc["machinepool"] != nil {
		o.machinePoolName = mc["machinepool"].(string)
	}
	switch i := mc["replicas"].(type) {
	case float64:
		o.replicas = int(i)
	case int:
		o.replicas = i
	}

	return nil
}

func (o *Options) validate() (err error) {

	if o.clusterName == "" {
		return fmt.Errorf("cluster name is missing")
	}

	if o.machinePoolName == "" {
		return fmt.Errorf("machinepool name is missing")
	}

	// // replicas defaults
	if o.replicas == 0 {
		return fmt.Errorf("replicas is missing")
	}

	return nil
}

func (o *Options) run() error {
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
	gvr := schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "machinepools"}
	dr := dynamicClient.Resource(gvr)
	mp, err := dr.Namespace(o.clusterName).Get(context.TODO(), o.machinePoolName, metav1.GetOptions{})
	if err != nil {
		return
	}
	if !o.CMFlags.DryRun {
		spec := mp.Object["spec"].(map[string]interface{})
		spec["replicas"] = o.replicas
		_, err = dr.Namespace(o.clusterName).Update(context.TODO(), mp, metav1.UpdateOptions{})
		if err != nil {
			return
		}
	}
	return nil
}
