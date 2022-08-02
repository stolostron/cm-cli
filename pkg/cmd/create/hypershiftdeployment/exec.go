// Copyright Contributors to the Open Cluster Management project
package hypershiftdeployment

import (
	"fmt"

	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/stolostron/applier/pkg/apply"
	"github.com/stolostron/cm-cli/pkg/cmd/create/hypershiftdeployment/scenario"
	"github.com/stolostron/cm-cli/pkg/helpers"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

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

const (
	azureLocationXPath = "managedCluster.infrastructure.platform.azure.location"
	awsRegionXPath     = "managedCluster.infrastructure.platform.aws.region"
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
	ok, err := helpers.NestedExists(o.values, "managedCluster")
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("managedCluster is missing")
	}

	if o.clusterName == "" {
		if o.clusterName, err = helpers.NestedString(o.values, "managedCluster.name"); err != nil {
			return err
		}
	}

	if err = helpers.SetNestedField(o.values, o.clusterName, "managedCluster.name"); err != nil {
		return err
	}

	if o.clusterNamespace == "" {
		if o.clusterNamespace, err = helpers.NestedString(o.values, "managedCluster.namespace"); err != nil {
			return err
		}
	}

	if err = helpers.SetNestedField(o.values, o.clusterNamespace, "managedCluster.namespace"); err != nil {
		return err
	}

	if o.hostingCluster == "" {
		if o.hostingCluster, err = helpers.NestedString(o.values, "managedCluster.hostingCluster"); err != nil {
			return err
		}
	}

	if err = helpers.SetNestedField(o.values, o.hostingCluster, "managedCluster.hostingCluster"); err != nil {
		return err
	}

	if o.hostingNamespace == "" {
		if o.hostingNamespace, err = helpers.NestedString(o.values, "managedCluster.hostingNamespace"); err != nil {
			return err
		}
	}

	if err = helpers.SetNestedField(o.values, o.hostingCluster, "managedCluster.hostingNamespace"); err != nil {
		return err
	}

	if o.cloudProviderSecretName == "" {
		if o.cloudProviderSecretName, err = helpers.NestedString(o.values,
			"managedCluster.infrastructure.cloudProvider.name"); err != nil {
			return err
		}
	}

	if err = helpers.SetNestedField(o.values, o.cloudProviderSecretName,
		"managedCluster.infrastructure.cloudProvider.name"); err != nil {
		return err
	}

	if o.region == "" {
		if ok, _ := helpers.NestedExists(o.values, awsRegionXPath); ok {
			if o.region, err = helpers.NestedString(o.values, awsRegionXPath); err != nil {
				return err
			}
		}
	}

	if o.region != "" {
		if err = helpers.SetNestedField(o.values, o.region, awsRegionXPath); err != nil {
			return err
		}
	}

	if o.location == "" {
		if ok, _ := helpers.NestedExists(o.values, azureLocationXPath); ok {
			if o.location, err = helpers.NestedString(o.values,
				azureLocationXPath); err != nil {
				return err
			}
		}
	}

	if o.location != "" {
		if err = helpers.SetNestedField(o.values,
			o.location,
			azureLocationXPath); err != nil {
			return err
		}
	}
	if len(o.region) != 0 && len(o.location) != 0 {
		return fmt.Errorf("only one plaform specification is allowed")
	}
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
	applierBuilder := apply.NewApplierBuilder()
	applier := applierBuilder.WithClient(kubeClient, apiextensionsClient, dynamicClient).Build()

	files := []string{
		"create/hd.yaml",
	}

	//overwrite nodepool name clusterName
	nodePools, ok, err := unstructured.NestedSlice(o.values, "managedCluster", "nodePools")
	if err != nil {
		return err
	}
	if ok {
		for i := range nodePools {
			nodePool := nodePools[i].(map[string]interface{})
			if err = unstructured.SetNestedField(nodePool, fmt.Sprintf("%s-%d", o.clusterName, i), "name"); err != nil {
				return err
			}
			spec, ok, err := unstructured.NestedMap(nodePool, "spec")
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("nodepool %s spec not found", nodePool["name"])
			}
			if err = unstructured.SetNestedField(spec, o.clusterName, "clusterName"); err != nil {
				return err
			}
			if err = unstructured.SetNestedField(nodePool, spec, "spec"); err != nil {
				return err
			}
		}
		if err := helpers.SetNestedField(o.values, nodePools, "managedCluster.nodePools"); err != nil {
			return err
		}
	}
	out, err := applier.ApplyCustomResources(reader, o.values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	return apply.WriteOutput(o.outputFile, output)
}
