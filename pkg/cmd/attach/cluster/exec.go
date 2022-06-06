// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
	workclientset "open-cluster-management.io/api/client/work/clientset/versioned"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/cmd/attach/cluster/scenario"
	"github.com/stolostron/cm-cli/pkg/helpers"
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
		reader := scenario.GetScenarioResourcesReader()
		o.values, err = helpers.ConvertReaderFileToValuesMap(valuesDefaultPath, reader)
		if err != nil {
			return err
		}
		if err = helpers.SetNestedField(o.values, o.clusterName, "managedCluster"); err != nil {
			return err
		}
	} else {
		//Read values
		o.values, err = helpers.ConvertValuesFileToValuesMap(o.valuesPath, "")
		if err != nil {
			return err
		}
	}

	ok, err := helpers.NestedExists(o.values, "managedCluster")
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("managedCluster is missing")
	}

	ok, err = helpers.NestedExists(o.values, "managedCluster.labels")
	if err != nil {
		return nil
	}
	if !ok {
		if err := helpers.SetNestedField(o.values, map[string]interface{}{
			"cloud":  "auto-detect",
			"vendor": "auto-detect",
		}, "managedCluster.labels"); err != nil {
			return err
		}
	}

	labels, _, err := unstructured.NestedMap(o.values, "managedCluster", "labels")
	if _, ok := labels["vendor"]; !ok {
		labels["vendor"] = "auto-detect"
	}

	if _, ok := labels["cloud"]; !ok {
		labels["cloud"] = "auto-detect"
	}

	if err = unstructured.SetNestedMap(o.values, labels, "managedCluster", "labels"); err != nil {
		return err
	}

	if o.clusterKubeConfig == "" && o.clusterKubeConfigContent == "" {
		kubeConfig, err := helpers.NestedString(o.values, "managedCluster.kubeConfig")
		if err == nil {
			o.clusterKubeConfig = kubeConfig
		}
	} else {
		if o.clusterKubeConfigContent == "" {
			b, err := ioutil.ReadFile(o.clusterKubeConfig)
			if err != nil {
				return err
			}
			o.clusterKubeConfig = string(b)
		} else {
			o.clusterKubeConfig = o.clusterKubeConfigContent
		}
	}

	if o.clusterKubeConfig != "" && o.clusterKubeContext != "" {
		config, err := clientcmd.Load([]byte(o.clusterKubeConfig))
		if err != nil {
			return err
		}
		if _, ok := config.Contexts[o.clusterKubeContext]; !ok {
			return fmt.Errorf("context %s doesn't exist in the provided kubeconfig", o.clusterKubeContext)
		}
		config.CurrentContext = o.clusterKubeContext
		b, err := clientcmd.Write(*config)
		if err != nil {
			return err
		}
		o.clusterKubeConfig = string(b)
	}

	if err := helpers.SetNestedField(o.values, o.clusterKubeConfig, "managedCluster.kubeConfig"); err != nil {
		return err
	}

	if o.clusterServer == "" {
		server, err := helpers.NestedString(o.values, "managedCluster.server")
		if err == nil {
			o.clusterServer = server
		}
	}

	if err := helpers.SetNestedField(o.values, o.clusterServer, "managedCluster.server"); err != nil {
		return err
	}

	if o.clusterToken == "" {
		token, err := helpers.NestedString(o.values, "managedCluster.token")
		if err == nil {
			o.clusterToken = token
		}
	}
	if err := helpers.SetNestedField(o.values, o.clusterToken, "managedCluster.token"); err != nil {
		return err
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

	if o.clusterName != "local-cluster" {
		if o.clusterKubeConfig != "" && (o.clusterToken != "" || o.clusterServer != "") {
			return fmt.Errorf("server/token and kubeConfig are mutually exclusif")
		}

		if (o.clusterToken == "" && o.clusterServer != "") ||
			(o.clusterToken != "" && o.clusterServer == "") {
			return fmt.Errorf("server or token is missing or should be removed")
		}

		if o.clusterKubeConfig != "" || o.clusterToken != "" {
			rhacmConstraint := ">=2.3.0"
			mceConstraint := ">=1.0.0"
			supported, platform, err := helpers.IsSupportedVersion(o.CMFlags, false, "", rhacmConstraint, mceConstraint)
			if err != nil {
				return err
			}
			if !supported {
				switch platform {
				case helpers.RHACM:
					return fmt.Errorf("auto-import is supported only on versions %s", rhacmConstraint)
				case helpers.MCE:
					return fmt.Errorf("auto-import is supported only on versions %s", mceConstraint)
				}
			}
		}

		//TODO must check if clusterDeployment CRD exists.
		gvr := schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterdeployments"}
		_, err := dynamicClient.Resource(gvr).Namespace(o.clusterName).Get(context.TODO(), o.clusterName, metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		} else {
			o.hiveScenario = true
		}

		if o.clusterKubeConfig == "" &&
			o.clusterToken == "" &&
			o.clusterServer == "" &&
			o.importFile == "" &&
			!o.hiveScenario {
			return fmt.Errorf("either kubeConfig or token/server or import-file must be provided")
		}
	}

	return nil
}

func (o *Options) run() (err error) {
	kubeClient, apiextensionsClient, dynamicClient, err := clusteradmhelpers.GetClients(o.CMFlags.KubectlFactory)
	if err != nil {
		return err
	}
	restConfig, err := o.CMFlags.KubectlFactory.ToRESTConfig()
	if err != nil {
		return err
	}
	clusterClient, err := clusterclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	workClient, err := workclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	return o.runWithClient(kubeClient, apiextensionsClient, dynamicClient, clusterClient, workClient)
}

func (o *Options) runWithClient(kubeClient kubernetes.Interface,
	apiextensionsClient apiextensionsclient.Interface,
	dynamicClient dynamic.Interface,
	clusterClient clusterclientset.Interface,
	workClient workclientset.Interface) (err error) {
	output := make([]string, 0)
	reader := scenario.GetScenarioResourcesReader()

	files := []string{
		"attach/hub/namespace.yaml",
	}

	if o.clusterKubeConfig != "" || o.clusterToken != "" {
		files = append(files, "attach/hub/managed_cluster_secret.yaml")
	}

	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.WithClient(kubeClient, apiextensionsClient, dynamicClient).Build()
	out, err := applier.ApplyDirectly(reader, o.values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	files = []string{
		"attach/hub/managed_cluster_cr.yaml",
	}

	if helpers.IsRHACM(o.CMFlags) {
		files = append(files, "attach/hub/klusterlet_addon_config_cr.yaml")
	}

	out, err = applier.ApplyCustomResources(reader, o.values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	if !o.hiveScenario &&
		o.importFile != "" &&
		!o.CMFlags.DryRun &&
		o.clusterName != "local-cluster" {
		time.Sleep(10 * time.Second)
		importSecret, err := kubeClient.CoreV1().
			Secrets(o.clusterName).
			Get(context.TODO(), o.clusterName+"-import", metav1.GetOptions{})
		if err != nil {
			return err
		}

		values := make(map[string]string)
		values["crds_yaml"] = string(importSecret.Data["crds.yaml"])
		values["import_yaml"] = string(importSecret.Data["import.yaml"])
		importFileContentCRD, err := applier.MustTemplateAsset(reader, values, "", "attach/managedcluster/import_crd.yaml")
		if err != nil {
			return err
		}
		importFileContentCRDFileName := fmt.Sprintf("%s_crd.yaml", o.importFile)
		importFileContentYAML, err := applier.MustTemplateAsset(reader, values, "", "attach/managedcluster/import_yaml.yaml")
		if err != nil {
			return err
		}
		importFileContentYAMLFileName := fmt.Sprintf("%s_yaml.yaml", o.importFile)

		err = ioutil.WriteFile(importFileContentCRDFileName, importFileContentCRD, 0600)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(importFileContentYAMLFileName, importFileContentYAML, 0600)
		if err != nil {
			return err
		}
		fmt.Printf("Execute this command on the managed cluster\nkubectl apply -f %s;sleep 10; kubectl apply -f %s\n",
			importFileContentCRDFileName,
			importFileContentYAMLFileName)
	}

	if !o.CMFlags.DryRun {
		if o.waitAgent || o.waitAddOns {
			return helpers.WaitKlusterlet(clusterClient, o.clusterName, o.timeout)
		}
		if o.waitAddOns {
			return helpers.WaitKlusterletAddons(workClient, o.clusterName, o.timeout)
		}
	}
	return clusteradmapply.WriteOutput(o.outputFile, output)
}
