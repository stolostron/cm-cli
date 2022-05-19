// Copyright Contributors to the Open Cluster Management project
package hypershiftdeployment

import (
	"fmt"

	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"

	"github.com/stolostron/cm-cli/pkg/cmd/create/hypershiftdeployment/scenario"
	"github.com/stolostron/cm-cli/pkg/helpers"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/spf13/cobra"
)

const (
	AWS       = "aws"
	AZURE     = "azure"
	GCP       = "gcp"
	OPENSTACK = "openstack"
	VSPHERE   = "vsphere"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	o.values, err = helpers.ConvertValuesFileToValuesMap(o.valuesPath, "")
	if err != nil {
		return err
	}

	if len(o.values) == 0 {
		return fmt.Errorf("values are missing")
	}

	if len(args) > 0 {
		o.clusterName = args[0]
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

	if o.clusterNamespace == "" {
		iname, ok := mc["namespace"]
		if !ok || iname == nil {
			return fmt.Errorf("cluster namespace is missing")
		}
		o.clusterNamespace = iname.(string)
		if len(o.clusterNamespace) == 0 {
			return fmt.Errorf("managedCluster.namespace not specified")
		}
	}

	mc["namespace"] = o.clusterNamespace

	if o.hostingCluster == "" {
		iname, ok := mc["hostingCluster"]
		if !ok || iname == nil {
			return fmt.Errorf("hostingCluster is missing")
		}
		o.hostingCluster = iname.(string)
		if len(o.hostingCluster) == 0 {
			return fmt.Errorf("managedCluster.hostingCluster not specified")
		}
	}

	mc["hostingCluster"] = o.hostingCluster

	if o.hostingNamespace == "" {
		iname, ok := mc["hostingNamespace"]
		if !ok || iname == nil {
			return fmt.Errorf("hostingNamespace is missing")
		}
		o.hostingNamespace = iname.(string)
		if len(o.hostingNamespace) == 0 {
			return fmt.Errorf("managedCluster.hostingNamespace not specified")
		}
	}

	mc["hostingNamespace"] = o.hostingNamespace

	if o.cloudProviderSecretName == "" {
		iname, ok := mc["cloudProviderSecretName"]
		if !ok || iname == nil {
			return fmt.Errorf("cloudProviderSecretName is missing")
		}
		o.cloudProviderSecretName = iname.(string)
		if len(o.cloudProviderSecretName) == 0 {
			return fmt.Errorf("managedCluster.cloudProviderSecretName not specified")
		}
	}

	mc["cloudProviderSecretName"] = o.cloudProviderSecretName

	if o.region == "" {
		iname, ok := mc["region"]
		if !ok || iname == nil {
			return fmt.Errorf("region is missing")
		}
		o.region = iname.(string)
		if len(o.region) == 0 {
			return fmt.Errorf("managedCluster.region not specified")
		}
	}

	mc["region"] = o.region

	return nil
}

func (o *Options) run() error {
	kubeClient, apiextensionsClient, dynamicClient, err := clusteradmhelpers.GetClients(o.CMFlags.KubectlFactory)
	if err != nil {
		return err
	}
	return o.runWithClient(kubeClient, apiextensionsClient, dynamicClient)
}

func (o *Options) runWithClient(kubeClient kubernetes.Interface,
	apiextensionsClient apiextensionsclient.Interface,
	dynamicClient dynamic.Interface) (err error) {
	output := make([]string, 0)
	reader := scenario.GetScenarioResourcesReader()
	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.WithClient(kubeClient, apiextensionsClient, dynamicClient).Build()

	files := []string{
		"create/hd.yaml",
	}

	out, err := applier.ApplyCustomResources(reader, o.values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	return clusteradmapply.WriteOutput(o.outputFile, output)
}
