// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	//Check if default values must be used
	if o.valuesPath == "" {
		if len(args) > 0 {
			o.clusterName = args[0]
		}
		if len(o.clusterName) == 0 {
			return fmt.Errorf("values or name are missing")
		}
		o.values = make(map[string]interface{})
		mc := make(map[string]interface{})
		mc["name"] = o.clusterName
		o.values["managedCluster"] = mc
	} else {
		//Read values
		o.values, err = helpers.ConvertValuesFileToValuesMap(o.valuesPath, "")
		if err != nil {
			return err
		}
	}

	if len(o.values) == 0 {
		return fmt.Errorf("values are missing")
	}

	return nil
}

func (o *Options) validate() (err error) {
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

func (o *Options) run() error {
	restConfig, err := o.CMFlags.KubectlFactory.ToRESTConfig()
	if err != nil {
		return err
	}
	clusterClient, err := clusterclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}
	return o.runWithClient(clusterClient, dynamicClient)
}

func (o *Options) runWithClient(clusterClient clusterclientset.Interface, dynamicClient dynamic.Interface) error {
	if !o.CMFlags.DryRun {
		err := clusterClient.ClusterV1().ManagedClusters().Delete(context.TODO(), o.clusterName, metav1.DeleteOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			fmt.Printf("managedcluster %s\n", err)
		}
		gvr := schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterdeployments"}
		err = dynamicClient.Resource(gvr).Namespace(o.clusterName).Delete(context.TODO(), o.clusterName, metav1.DeleteOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			fmt.Printf("clusterdeployment %s\n", err)
		}
	}
	return nil
}
