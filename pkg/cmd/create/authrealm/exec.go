// Copyright Contributors to the Open Cluster Management project
package authrealm

import (
	"context"
	"fmt"

	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"

	"github.com/stolostron/applier/pkg/apply"
	"github.com/stolostron/cm-cli/pkg/cmd/create/authrealm/scenario"
	"github.com/stolostron/cm-cli/pkg/helpers"
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

	if len(o.placement) > 0 {
		authRealm["placement"] = o.placement
	}

	if len(o.managedClusterSet) > 0 {
		authRealm["managedClusterSet"] = o.managedClusterSet
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
	if _, ok := authRealm["identityProviders"]; !ok {
		return fmt.Errorf("identityProviders is missing")
	}

	if v, ok := authRealm["placement"]; ok && v != nil {
		if _, err = dynamicClient.Resource(helpers.GvrPLC).
			Namespace(authRealm["namespace"].(string)).
			Get(context.TODO(), v.(string), metav1.GetOptions{}); err != nil {
			return err
		}
	}
	if v, ok := authRealm["managedClusterSet"]; ok && v != nil {
		if _, err = dynamicClient.Resource(helpers.GvrMCS).
			Get(context.TODO(), v.(string), metav1.GetOptions{}); err != nil {
			return err
		}
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
	applierBuilder := apply.NewApplierBuilder()
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
	if placement, ok := authRealm["placement"]; !ok || placement == nil {
		file := "create/hub/placement.yaml"
		out, err := applier.ApplyCustomResource(reader, o.values, o.CMFlags.DryRun, "", file)
		if err != nil {
			return err
		}
		output = append(output, out)
		authRealm["placement"] = authRealm["name"].(string) + "-placement"
	}

	//create clusterset if not provided
	if managedClusterSet, ok := authRealm["managedClusterSet"]; !ok || managedClusterSet == nil {
		file := "create/hub/managedclusterset.yaml"
		out, err := applier.ApplyCustomResource(reader, o.values, o.CMFlags.DryRun, "", file)
		if err != nil {
			return err
		}
		output = append(output, out)
		authRealm["managedClusterSet"] = authRealm["name"].(string) + "-clusterset"
	}

	//Create binding
	{
		file = "create/hub/managedclusterset_binding.yaml"
		out, err := applier.ApplyCustomResource(reader, o.values, o.CMFlags.DryRun, "", file)
		if err != nil {
			return err
		}
		output = append(output, out)
		authRealm["managedClusterSetBinding"] = authRealm["managedClusterSet"]
	}

	//Create secrets
	for _, iidp := range authRealm["identityProviders"].([]interface{}) {
		idp := iidp.(map[string]interface{})
		o.values["idp"] = idp
		switch idp["type"].(string) {
		case "GitHub":
			file := "create/hub/github_secret.yaml"
			out, err := applier.ApplyDirectly(reader, o.values, o.CMFlags.DryRun, "", file)
			if err != nil {
				return err
			}
			igithub := idp["github"]
			github := igithub.(map[string]interface{})
			github["clientSecret"] = map[string]string{"name": idp["name"].(string) + "-secret"}
			output = append(output, out...)
		case "LDAP":
			file := "create/hub/ldap_secret.yaml"
			out, err := applier.ApplyDirectly(reader, o.values, o.CMFlags.DryRun, "", file)
			if err != nil {
				return err
			}
			ildap := idp["ldap"]
			ldap := ildap.(map[string]interface{})
			ldap["bindPassword"] = map[string]string{"name": idp["name"].(string) + "-secret"}
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

	return apply.WriteOutput(o.outputFile, output)
}
