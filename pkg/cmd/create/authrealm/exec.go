// Copyright Contributors to the Open Cluster Management project
package authrealm

import (
	"context"
	"fmt"

	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"

	"github.com/open-cluster-management/cm-cli/pkg/cmd/create/authrealm/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	o.values, err = helpers.ConvertValuesFileToValuesMap(o.valuesPath, "")
	if err != nil {
		return err
	}

	if len(o.values) == 0 {
		return fmt.Errorf("values are missing")
	}

	if len(o.name) == 0 && len(args) > 0 {
		o.name = args[0]
	}
	iauthRealm := o.values["authRealm"]
	authRealm := iauthRealm.(map[string]interface{})
	if len(o.name) > 0 {
		authRealm["name"] = o.name
	}

	if len(o.namespace) > 0 {
		authRealm["namespace"] = o.namespace
	}

	if len(o.typeName) > 0 {
		authRealm["type"] = o.typeName
	}

	if len(o.routeSubDomain) > 0 {
		authRealm["routeSubDomain"] = o.routeSubDomain
	}

	if len(o.placementName) > 0 {
		authRealm["placementName"] = o.placementName
	}

	if len(o.clusterSetName) > 0 {
		authRealm["clusterSetName"] = o.clusterSetName
	}

	if len(o.clusterSetBindingName) > 0 {
		authRealm["clusterSetBindingName"] = o.clusterSetBindingName
	}

	return nil

}

func (o *Options) validate() (err error) {
	_, apiExtensionClient, dynamicClient, err := clusteradmhelpers.GetClients(o.CMFlags.KubectlFactory)
	if err != nil {
		return err
	}
	if _, err := apiExtensionClient.
		ApiextensionsV1().
		CustomResourceDefinitions().
		Get(context.TODO(), "authrealms.identityconfig.identitatem.io", metav1.GetOptions{}); err != nil {
		return fmt.Errorf("identitatem not installed")
	}
	iauthRealm := o.values["authRealm"]
	authRealm := iauthRealm.(map[string]interface{})
	if _, ok := authRealm["name"]; !ok {
		return fmt.Errorf("name is missing")
	}
	if _, ok := authRealm["namespace"]; !ok {
		return fmt.Errorf("namespace is missing")
	}
	if _, ok := authRealm["type"]; !ok {
		return fmt.Errorf("type is missing")
	}
	if v, ok := authRealm["type"]; ok && v != "dex" {
		return fmt.Errorf("type possible values is dex")
	}
	if _, ok := authRealm["routeSubDomain"]; !ok {
		return fmt.Errorf("routeSubDomain is missing")
	}

	if v, ok := authRealm["placementName"]; ok && v != nil {
		if _, err = dynamicClient.Resource(helpers.GvrPLC).
			Namespace(authRealm["namespace"].(string)).
			Get(context.TODO(), v.(string), metav1.GetOptions{}); err != nil {
			return err
		}
	}
	if v, ok := authRealm["managedClusterSetName"]; ok && v != nil {
		if _, err = dynamicClient.Resource(helpers.GvrMCS).
			Namespace(authRealm["namespace"].(string)).
			Get(context.TODO(), v.(string), metav1.GetOptions{}); err != nil {
			return err
		}
	}
	if v, ok := authRealm["managedClusterSetBindingName"]; ok && v != nil {
		if _, err = dynamicClient.Resource(helpers.GvrMCSB).
			Namespace(authRealm["namespace"].(string)).
			Get(context.TODO(), v.(string), metav1.GetOptions{}); err != nil {
			return err
		}
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
	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.WithClient(kubeClient, apiextensionsClient, dynamicClient).Build()

	iauthRealm := o.values["authRealm"]
	authRealm := iauthRealm.(map[string]interface{})

	//create ns
	file := "create/hub/namespace.yaml"
	out, err := applier.ApplyDirectly(reader, o.values, o.CMFlags.DryRun, "", file)
	if err != nil {
		return err
	}
	output = append(output, out...)

	//create placement if not provided
	if placementName, ok := authRealm["placementName"]; !ok || placementName == nil {
		file := "create/hub/placement.yaml"
		out, err := applier.ApplyCustomResource(reader, o.values, o.CMFlags.DryRun, "", file)
		if err != nil {
			return err
		}
		output = append(output, out)
		authRealm["placementName"] = authRealm["name"].(string) + "-placement"
	}

	//create clusterset if not provided
	if managedClusterSetName, ok := authRealm["managedClusterSetName"]; !ok || managedClusterSetName == nil {
		file := "create/hub/managedclusterset.yaml"
		out, err := applier.ApplyCustomResource(reader, o.values, o.CMFlags.DryRun, "", file)
		if err != nil {
			return err
		}
		output = append(output, out)
		authRealm["managedClusterSetName"] = authRealm["name"].(string) + "-clusterset"
	}

	//create clustersetbinding if not provided
	if managedClusterSetBindingName, ok := authRealm["managedClusterSetBindingName"]; !ok || managedClusterSetBindingName == nil {
		file := "create/hub/managedclusterset_binding.yaml"
		out, err := applier.ApplyCustomResource(reader, o.values, o.CMFlags.DryRun, "", file)
		if err != nil {
			return err
		}
		output = append(output, out)
		authRealm["managedClusterSetBindingName"] = authRealm["name"].(string) + "-clusterset"
	}

	//Create secrets
	for _, iidp := range authRealm["identityProviders"].([]interface{}) {
		idp := iidp.(map[string]interface{})
		o.values["idp"] = idp
		switch idp["type"].(string) {
		case "GitHub":
			file := "create/hub/secret.yaml"
			out, err := applier.ApplyDirectly(reader, o.values, o.CMFlags.DryRun, "", file)
			if err != nil {
				return err
			}
			igithub := idp["github"]
			github := igithub.(map[string]interface{})
			github["clientSecret"] = map[string]string{"name": authRealm["name"].(string) + "-client-secret"}
			output = append(output, out...)
		default:
			return fmt.Errorf("unsupported idp type %s", idp["type"].(string))
		}
	}

	file = "create/hub/authrealm.yaml"
	out, err = applier.ApplyCustomResources(reader, o.values, o.CMFlags.DryRun, "", file)
	if err != nil {
		return err
	}
	output = append(output, out...)

	return clusteradmapply.WriteOutput(o.outputFile, output)
}
