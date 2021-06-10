// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"context"
	"fmt"
	"path/filepath"

	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"

	"github.com/open-cluster-management/cm-cli/pkg/cmd/create/cluster/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ghodss/yaml"

	"k8s.io/client-go/discovery"
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

	return nil
}

func (o *Options) validate() (err error) {
	imc, ok := o.values["managedCluster"]
	if !ok || imc == nil {
		return fmt.Errorf("managedCluster is missing")
	}
	mc := imc.(map[string]interface{})
	icloud, ok := mc["cloud"]
	if !ok || icloud == nil {
		return fmt.Errorf("cloud type is missing")
	}
	cloud := icloud.(string)
	if cloud != AWS && cloud != AZURE && cloud != GCP && cloud != OPENSTACK && cloud != VSPHERE {
		return fmt.Errorf("supported cloud type are (%s, %s, %s, %s, %s) and got %s", AWS, AZURE, GCP, OPENSTACK, VSPHERE, cloud)
	}
	o.cloud = cloud

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
	kubeClient, err := o.CMFlags.KubectlFactory.KubernetesClientSet()
	if err != nil {
		return err
	}
	dynamicClient, err := o.CMFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}
	restConfig, err := o.CMFlags.KubectlFactory.ToRESTConfig()
	if err != nil {
		return err
	}
	apiextensionsClient, err := apiextensionsclient.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	discoveryClient, err := o.CMFlags.KubectlFactory.ToDiscoveryClient()
	if err != nil {
		return err
	}
	return o.runWithClient(kubeClient, dynamicClient, apiextensionsClient, discoveryClient)
}

func (o *Options) runWithClient(kubeClient kubernetes.Interface,
	dynamicClient dynamic.Interface,
	apiextensionsClient apiextensionsclient.Interface,
	discoveryClient discovery.DiscoveryInterface) (err error) {
	output := make([]string, 0)
	pullSecret, err := kubeClient.CoreV1().Secrets("openshift-config").Get(
		context.TODO(),
		"pull-secret",
		metav1.GetOptions{})
	if err != nil {
		return err
	}

	ps, err := yaml.Marshal(pullSecret)
	if err != nil {
		return err
	}

	valueps := make(map[string]interface{})
	err = yaml.Unmarshal(ps, &valueps)
	if err != nil {
		return err
	}

	o.values["pullSecret"] = valueps

	reader := scenario.GetScenarioResourcesReader()
	installConfig, err := clusteradmapply.MustTempalteAsset(reader, o.values, "", filepath.Join(scenarioDirectory, "hub", o.cloud, "install_config.yaml"))
	if err != nil {
		return err
	}

	valueic := make(map[string]interface{})
	err = yaml.Unmarshal(installConfig, &valueic)
	if err != nil {
		return err
	}

	files := []string{
		"create/hub/common/namespace.yaml",
	}

	clientHolder := resourceapply.NewClientHolder().
		WithAPIExtensionsClient(apiextensionsClient).
		WithKubernetes(kubeClient).
		WithDynamicClient(dynamicClient)

	out, err := clusteradmapply.ApplyDirectly(clientHolder, reader, o.values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)
	o.values["installConfig"] = valueic

	files = []string{

		"create/hub/common/creds_secret_cr.yaml",
		"create/hub/common/install_config_secret_cr.yaml",
		"create/hub/common/klusterlet_addon_config_cr.yaml",
		"create/hub/common/machinepool_cr.yaml",
		"create/hub/common/pull_secret_cr.yaml",
		"create/hub/common/ssh_private_key_secret_cr.yaml",
		"create/hub/common/vsphere_ca_cert_secret_cr.yaml",
		"create/hub/common/clusterimageset_cr.yaml",
	}

	files = append(files,
		"create/hub/common/cluster_deployment_cr.yaml",
		"create/hub/common/managed_cluster_cr.yaml")

	out, err = clusteradmapply.ApplyCustomResouces(dynamicClient, discoveryClient, reader, o.values, o.CMFlags.DryRun, "create/hub/common/_helpers.tpl", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)
	return clusteradmapply.WriteOutput(o.outputFile, output)
}
