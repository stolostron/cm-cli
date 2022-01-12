// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"context"
	"fmt"
	"path/filepath"

	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
	workclientset "open-cluster-management.io/api/client/work/clientset/versioned"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"

	attachscenario "github.com/stolostron/cm-cli/pkg/cmd/attach/cluster/scenario"
	"github.com/stolostron/cm-cli/pkg/cmd/create/cluster/scenario"
	"github.com/stolostron/cm-cli/pkg/helpers"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ghodss/yaml"

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
	icloud, ok := mc["cloud"]
	if !ok || icloud == nil {
		return fmt.Errorf("cloud type is missing")
	}
	cloud := icloud.(string)
	if cloud != AWS && cloud != AZURE && cloud != GCP && cloud != OPENSTACK && cloud != VSPHERE {
		return fmt.Errorf("supported cloud type are (%s, %s, %s, %s, %s) and got %s", AWS, AZURE, GCP, OPENSTACK, VSPHERE, cloud)
	}
	o.cloud = cloud

	_, ocpImageOk := mc["ocpImage"]
	_, imageSetRef := mc["imageSetRef"]
	if ocpImageOk && imageSetRef {
		return fmt.Errorf("ocpImage and imageSetRef are mutually exclusive")
	}

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

	if o.clusterSetName == "" {
		if iclusterSetName, ok := mc["clusterSetName"]; ok {
			o.clusterSetName = iclusterSetName.(string)
		}
	}

	mc["clusterSetName"] = o.clusterSetName

	return nil
}

func (o *Options) run() error {
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
	attachreader := attachscenario.GetScenarioResourcesReader()
	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.WithClient(kubeClient, apiextensionsClient, dynamicClient).Build()

	installConfig, err := applier.MustTempalteAsset(reader, o.values, "", filepath.Join(scenarioDirectory, "hub", o.cloud, "install_config.yaml"))
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

	out, err := applier.ApplyDirectly(reader, o.values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)
	o.values["installConfig"] = valueic

	files = []string{
		"create/hub/common/creds_secret_cr.yaml",
		"create/hub/common/install_config_secret_cr.yaml",
		"create/hub/common/machinepool_cr.yaml",
		"create/hub/common/pull_secret_cr.yaml",
		"create/hub/common/ssh_private_key_secret_cr.yaml",
		"create/hub/common/vsphere_ca_cert_secret_cr.yaml",
	}

	imc := o.values["managedCluster"]
	mc := imc.(map[string]interface{})

	if _, ok := mc["imageSetRef"]; !ok {
		files = append(files, "create/hub/common/clusterimageset_cr.yaml")
	}

	files = append(files, "create/hub/common/cluster_deployment_cr.yaml")
	out, err = applier.ApplyCustomResources(reader, o.values, o.CMFlags.DryRun, "create/hub/common/_helpers.tpl", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	files = []string{
		"attach/hub/managed_cluster_cr.yaml",
		"attach/hub/klusterlet_addon_config_cr.yaml",
	}
	out, err = applier.ApplyCustomResources(attachreader, o.values, o.CMFlags.DryRun, "", files...)
	if err != nil {
		return err
	}

	output = append(output, out...)
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
