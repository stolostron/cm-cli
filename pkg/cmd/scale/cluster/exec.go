// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"context"
	"fmt"
	"path/filepath"

	appliercmd "github.com/open-cluster-management/applier/pkg/applier/cmd"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"

	//libgoclient "github.com/open-cluster-management/library-go/pkg/client"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"

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

	//if o == nil {
	//	o = map[string]interface{}{}
	//}

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

	if o.machinePoolName == "" {
		return fmt.Errorf("machinepool is missing")
	}

	// replicas defaults

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

	reader := scenario.GetApplierScenarioResourcesReader()
	mp := unstructured.Unstructured{}
	//mp.SetKind("MachinePool")
	//mp.SetAPIVersion("hive.openshift.io/v1")

	mp.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "hive.openshift.io",
		Kind:    "MachinePool",
		Version: "v1",
	})

	klog.Info("get MachinePool entry")
	err := client.Get(context.TODO(),
		crclient.ObjectKey{
			Name:      o.machinePoolName,
			Namespace: o.clusterName,
		}, &mp)
	klog.Info("get MachinePool entry - return.  %v", mp)

	if err == nil {

		//	gvrMachinePool = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "machinepool"}
		//	clientHubDynamic, err = libgoclient.NewDefaultKubeClientDynamic("")
		//	libgoclient.

		// add place holder
		//patchString := `{"spec": {"replicas": 3 }}`
		patchInBytes := []byte(fmt.Sprintf(`{"spec":{replicas":"%d"}}`, o.replicas))

		//		_, err := clientHubDynamic.Resource(gvrMachinePool).Namespace("").Patch(context.TODO(), name, types.MergePatchType, []byte(patchString), metav1.PatchOptions{}, "status")
		//_, err := crclient.Patch(context.TODO(), mp, types.MergePatchType, []byte(patchString), metav1.PatchOptions{})
		//_, err := crclient.Patch(context.TODO(), mp, crclient.RawPatch(types.MergePatchType, patch))
		//_, err := crclient.Patch(context.TODO(), mp, types.JSONPatchType, patchInBytes, v1.PatchOptions{})
		//_, err := crclient.Patch(context.TODO(), mp, crclient.RawPatch(types.MergePatchType, patchInBytes))

		klog.V(1).Info("Change replicas to %d", o.replicas)

		//crclient.RawPatch(types.StrategicMergePatchType, patchInBytes)
		//crclient.RawPatch(types.MergePatchType, patchInBytes)
		//crclient.RawPatch(types.ApplyPatchType, patchInBytes)
		crclient.RawPatch(types.JSONPatchType, patchInBytes)

		//  patchStringValue specifies a json patch operation for a string.
		//	type patchStringValue struct {
		//		Op    string `json:"op"`
		//		Path  string `json:"path"`
		//		Value string `json:"value"`
		//	}

		//	patch := []patchStringValue{{
		//		Op:    "replace",
		//		Path:  "/spec/replicas",
		//		Value: strconv.Itoa(o.replicas),
		//	}}
		//patchInBytes, _ := json.Marshal(patch)
		//klog.V(2).Info(" > Patching secret " + secretName + " in namespace " + clusterName)
		//_, err = crclient.Patch(context.TODO(), mp, types.JSONPatchType, patchInBytes, v1.PatchOptions{})

		//_, err = crclient.Patch(context.TODO(), mp, types.JSONPatchType, patchInBytes, metav1.PatchOptions{})

		if err != nil {
			panic(fmt.Errorf("Failed to scale replicas: %v", err))
		}

		//Another try

		//		//spec.replicas
		//		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		//	// Retrieve the latest version of Deployment before attempting update
		//	// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		//	result, getErr := client.Resource(deploymentRes).Namespace(namespace).Get("demo-deployment", metav1.GetOptions{})
		//	if getErr != nil {
		//		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		//	}
		//
		//	// update replicas to 1
		//	if err := unstructured.SetNestedField(result.Object, int64(1), "spec", "replicas"); err != nil {
		//		panic(fmt.Errorf("Failed to set replica value: %v", err))
		//	}

	} else {
		return fmt.Errorf("MachinePool not found! %v", err)
	}

	applyOptions := &appliercmd.Options{
		OutFile:     o.applierScenariosOptions.OutFile,
		ConfigFlags: o.applierScenariosOptions.ConfigFlags,

		Delete:    false,
		Timeout:   o.applierScenariosOptions.Timeout,
		Force:     o.applierScenariosOptions.Force,
		Silent:    o.applierScenariosOptions.Silent,
		IOStreams: o.applierScenariosOptions.IOStreams,
	}

	return applyOptions.ApplyWithValues(client, reader,
		filepath.Join(scenarioDirectory, "hub", "common"),
		o.values)
}
