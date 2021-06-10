// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"fmt"

	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"

	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	//Check if default values must be used
	if o.valuesPath == "" {
		if o.clusterName != "" {
			o.values = make(map[string]interface{})
			mc := make(map[string]interface{})
			mc["name"] = o.clusterName
			o.values["managedCluster"] = mc
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

	if len(o.values) == 0 {
		return fmt.Errorf("values are missing")
	}

	return nil
}

func (o *Options) validate() error {
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
	return o.runWithClient(clusterClient)
}

func (o *Options) runWithClient(clusterClient clusterclientset.Interface) error {
	if !o.CMFlags.DryRun {
		return clusterClient.ClusterV1().ManagedClusters().Delete(context.TODO(), o.clusterName, metav1.DeleteOptions{})
	}
	return nil
}
