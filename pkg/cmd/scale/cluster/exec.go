// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"context"
	"fmt"

	appliercmd "github.com/open-cluster-management/applier/pkg/applier/cmd"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/open-cluster-management/cm-cli/pkg/cmd/scale/cluster/scenario"

	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		return nil
	}
	o.values, err = appliercmd.ConvertValuesFileToValuesMap(o.applierScenariosOptions.ValuesPath, "")
	if err != nil {
		return err
	}

	return nil
}

func (o *Options) validate() (err error) {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		return nil
	}

	if o.clusterName == "" {
		return fmt.Errorf("cluster name is missing")
	}

	if o.machinePoolName == "" {
		return fmt.Errorf("machinepool is missing")
	}

	// // replicas defaults
	if o.replicas == 0 {
		return fmt.Errorf("replicas is missing")
	}

	return nil
}

func (o *Options) run() error {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		reader := scenario.GetApplierScenarioResourcesReader()
		return reader.ExtractAssets(scenarioDirectory, o.applierScenariosOptions.OutTemplatesDir)
	}
	client, err := helpers.GetControllerRuntimeClientFromFlags(o.applierScenariosOptions.ConfigFlags)
	if err != nil {
		return err
	}
	return o.runWithClient(client)
}

func (o *Options) runWithClient(client crclient.Client) error {
	mp := &unstructured.Unstructured{}
	mp.SetKind("MachinePool")
	mp.SetAPIVersion("hive.openshift.io/v1")
	err := client.Get(context.TODO(),
		crclient.ObjectKey{
			Name:      o.machinePoolName,
			Namespace: o.clusterName}, mp)
	if err != nil {
		return err
	}
	patch := crclient.MergeFrom(mp.DeepCopyObject())
	spec := mp.Object["spec"].(map[string]interface{})
	spec["replicas"] = o.replicas
	return client.Patch(context.TODO(), mp, patch)
}
