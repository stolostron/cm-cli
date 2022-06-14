// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/cmd/get/config/cluster/scenario"
	"github.com/stolostron/cm-cli/pkg/helpers"
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

	//Check if platform is supported.
	var secretName string
	switch {
	case cd.Spec.Platform.AWS != nil:
		secretName = cd.Spec.Platform.AWS.CredentialsSecretRef.Name
	case cd.Spec.Platform.Azure != nil && o.CMFlags.Beta:
		secretName = cd.Spec.Platform.Azure.CredentialsSecretRef.Name
	case cd.Spec.Platform.GCP != nil && o.CMFlags.Beta:
		secretName = cd.Spec.Platform.GCP.CredentialsSecretRef.Name
	case cd.Spec.Platform.VSphere != nil && o.CMFlags.Beta:
		secretName = cd.Spec.Platform.VSphere.CredentialsSecretRef.Name
	case cd.Spec.Platform.OpenStack != nil && o.CMFlags.Beta:
		secretName = cd.Spec.Platform.OpenStack.CredentialsSecretRef.Name
	default:
		return fmt.Errorf("unsupported platform %v", cd.Spec.Platform)
	}

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
	if o.withoutCredentials {
		values["awsAccessKeyID"] = "your_aws_access_key_id"
		values["awsSecretAccessKey"] = "your_aws_secret_access_key"
		values["osServicePrincipalJson"] = map[string]interface{}{
			"clientID":       "your_clientID",
			"clientSecret":   "your_clientSecret",
			"tenantID":       "your_tenantID",
			"subscriptionID": "your_subscriptionID",
		}
		values["osServiceAccountJson"] = "your_osServiceAccountJson"
		values["vsphere_username"] = "your_username"
		values["vsphere_password"] = "your_password"
		values["openstack_cloudsYaml"] = "your_cloudsYaml"
	} else {
		cred, err := kubeClient.CoreV1().Secrets(o.ClusterName).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		switch {
		case cd.Spec.Platform.AWS != nil:
			values["awsAccessKeyID"] = string(cred.Data["aws_access_key_id"])
			values["awsSecretAccessKey"] = string(cred.Data["aws_secret_access_key"])
		case cd.Spec.Platform.Azure != nil:
			osServicePrincipal := cred.Data["osServicePrincipal.json"]
			osServicePrincipalMap := make(map[string]interface{})
			err = json.Unmarshal(osServicePrincipal, &osServicePrincipalMap)
			if err != nil {
				return err
			}
			values["osServicePrincipalJson"] = osServicePrincipalMap
		case cd.Spec.Platform.GCP != nil:
			values["osServiceAccountJson"] = string(cred.Data["osServiceAccount.json"])
		case cd.Spec.Platform.VSphere != nil:
			values["vsphere_username"] = string(cred.Data["username"])
			values["vsphere_password"] = string(cred.Data["password"])
			cert, err := kubeClient.CoreV1().Secrets(o.ClusterName).Get(context.TODO(), cd.Name+"-vsphere-certs", metav1.GetOptions{})
			if err != nil {
				return err
			}
			//Not sure if I have to decode or not as the secret template contains a encode statement.
			values["vpshere_cert"] = string(cert.Data[".cacert"])
		case cd.Spec.Platform.OpenStack != nil:
			values["openstack_cloudsYaml"] = string(cred.Data["clouds.yaml"])
		}
	}

	//Get clusterimageset
	klog.V(5).Infof("ImageSetRef:%s", cd.Spec.Provisioning.ImageSetRef.Name)
	cisu, err := dynamicClient.Resource(helpers.GvrCIS).Get(context.TODO(), cd.Spec.Provisioning.ImageSetRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	values["clusterImageSet"] = cisu.Object

	klog.V(5).Infof("%v\n", values)
	applierBuilder := clusteradmapply.NewApplierBuilder()
	applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient).Build()
	b, err := applier.MustTemplateAsset(reader, values, "", "config/config.yaml")
	if err != nil {
		return err
	}
	if len(o.outputFile) != 0 {
		return ioutil.WriteFile(o.outputFile, b, 0600)
	}
	fmt.Printf("%s\n", string(b))
	return nil
}
