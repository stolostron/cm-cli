// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/get/config/cluster/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/spf13/cobra"
	apiextensionsClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	clusteradmhelpers "open-cluster-management.io/clusteradm/pkg/helpers"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 0 {
		o.ClusterName = args[0]
	}
	return nil
}

func (o *Options) validate() error {
	if len(o.ClusterName) == 0 {
		return fmt.Errorf("clustername is missing")
	}
	return nil
}

func (o *Options) run() (err error) {
	kubeClient, apiextensionsClient, dynamicClient, err := clusteradmhelpers.GetClients(o.CMFlags.KubectlFactory)
	if err != nil {
		return err
	}
	return o.runWithClient(kubeClient, apiextensionsClient, dynamicClient)
}

func (o *Options) runWithClient(kubeClient kubernetes.Interface,
	apiExtensionsClient apiextensionsClient.Interface,
	dynamicClient dynamic.Interface) (err error) {
	reader := scenario.GetScenarioResourcesReader()
	values := make(map[string]interface{})

	//Get clusterDeployment
	cdu, err := dynamicClient.Resource(helpers.GvrCD).Namespace(o.ClusterName).Get(context.TODO(), o.ClusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cd := &hivev1.ClusterDeployment{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd)
	if err != nil {
		return err
	}
	values["clusterDeployment"] = cdu.Object

	//Get install-config
	ic, err := kubeClient.CoreV1().Secrets(o.ClusterName).Get(context.TODO(), cd.Spec.Provisioning.InstallConfigSecretRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	j, err := yaml.YAMLToJSON(ic.Data["install-config.yaml"])
	if err != nil {
		return err
	}

	u := &unstructured.Unstructured{}
	_, _, err = unstructured.UnstructuredJSONScheme.Decode(j, nil, u)
	if err != nil {
		if !runtime.IsMissingKind(err) {
			return err
		}
	}
	values["installConfig"] = u.Object

	//Get sshPrivateKey
	pk, err := kubeClient.CoreV1().Secrets(o.ClusterName).Get(context.TODO(), cd.Spec.Provisioning.SSHPrivateKeySecretRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	values["sshPrivateKey"] = string(pk.Data["ssh-privatekey"])

	//Get credentials
	switch {
	case cd.Spec.Platform.AWS != nil:
		cred, err := kubeClient.CoreV1().Secrets(o.ClusterName).Get(context.TODO(), cd.Spec.Platform.AWS.CredentialsSecretRef.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		values["awsAccessKeyID"] = string(cred.Data["aws_access_key_id"])
		values["awsSecretAccessKey"] = string(cred.Data["aws_secret_access_key"])
	default:
		return fmt.Errorf("unsupported platform")
	}

	//Get clusterimageset
	klog.V(5).Infof("ImageSetRef:%s", cd.Spec.Provisioning.ImageSetRef.Name)
	cisu, err := dynamicClient.Resource(helpers.GvrCIS).Get(context.TODO(), cd.Spec.Provisioning.ImageSetRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	values["clusterImageSet"] = cisu.Object

	klog.V(5).Infof("%v\n", values)
	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient)
	b, err := applier.MustTempalteAsset(reader, values, "", "config/config.yaml")
	if err != nil {
		return err
	}
	if len(o.outputFile) != 0 {
		return ioutil.WriteFile(o.outputFile, b, 0600)
	}
	fmt.Printf("%s\n", string(b))
	return nil
}
