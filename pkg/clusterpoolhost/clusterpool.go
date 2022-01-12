// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	printclusterpoolv1alpha1 "github.com/stolostron/cm-cli/api/cm-cli/v1alpha1"
	"github.com/stolostron/cm-cli/pkg/clusterpoolhost/scenario"
	"github.com/stolostron/cm-cli/pkg/helpers"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"
)

func (cph *ClusterPoolHost) SizeClusterPool(clusterPoolName string, size int32, dryRun bool) error {
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}
	cpu, err := dynamicClient.Resource(helpers.GvrCP).Namespace(cph.Namespace).Get(context.TODO(), clusterPoolName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cp := &hivev1.ClusterPool{}
	if runtime.DefaultUnstructuredConverter.FromUnstructured(cpu.UnstructuredContent(), cp); err != nil {
		return err
	}
	if !dryRun {
		cp.Spec.Size = size
		cpu.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(cp)
		if runtime.DefaultUnstructuredConverter.FromUnstructured(cpu.UnstructuredContent(), cp); err != nil {
			return err
		}
		if _, err = dynamicClient.Resource(helpers.GvrCP).Namespace(cph.Namespace).Update(context.TODO(), cpu, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (cph *ClusterPoolHost) CreateClusterPool(clusterPoolName, cloud string, values map[string]interface{}, dryRun bool, outputFile string) error {
	values["namespace"] = cph.Namespace

	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	apiExtensionsClient, err := apiextensionsclient.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	output := make([]string, 0)

	reader := scenario.GetScenarioResourcesReader()
	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient).Build()

	installConfig, err := applier.MustTempalteAsset(reader, values, "", filepath.Join("create", "clusterpool", cloud, "install_config.yaml"))
	if err != nil {
		return err
	}

	valueic := make(map[string]interface{})
	err = yaml.Unmarshal(installConfig, &valueic)
	if err != nil {
		return err
	}

	files := []string{
		"create/clusterpool/common/namespace.yaml",
	}

	out, err := applier.ApplyDirectly(reader, values, dryRun, "", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)
	values["installConfig"] = valueic

	files = []string{
		"create/clusterpool/common/creds_secret_cr.yaml",
		"create/clusterpool/common/install_config_secret_cr.yaml",
		"create/clusterpool/common/pull_secret_cr.yaml",
	}

	cpi := values["clusterPool"]
	cp := cpi.(map[string]interface{})
	if _, ok := cp["imageSetRef"]; ok {
		files = append(files,
			"create/clusterpool/common/clusterpool_cr.yaml")
	} else {
		files = append(files, "create/clusterpool/common/clusterimageset_cr.yaml",
			"create/clusterpool/common/clusterpool_cr.yaml")
	}
	out, err = applier.ApplyCustomResources(reader, values, dryRun, "create/clusterpool/common/_helpers.tpl", files...)
	if err != nil {
		return err
	}
	output = append(output, out...)

	return clusteradmapply.WriteOutput(outputFile, output)
}

func (cph *ClusterPoolHost) DeleteClusterPools(clusterPoolNames string, dryRun bool, outputFile string) error {
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	if !dryRun {
		for _, cpn := range strings.Split(clusterPoolNames, ",") {
			clusterPoolName := strings.TrimSpace(cpn)
			err = dynamicClient.Resource(helpers.GvrCP).Namespace(cph.Namespace).Delete(context.TODO(), clusterPoolName, metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (cph *ClusterPoolHost) GetClusterPools(showCphName, dryRun bool) (*hivev1.ClusterPoolList, error) {
	clusterPools := &hivev1.ClusterPoolList{}
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return clusterPools, err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return clusterPools, err
	}

	l, err := dynamicClient.Resource(helpers.GvrCP).Namespace(cph.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return clusterPools, err
	}
	for _, cpu := range l.Items {
		cp := &hivev1.ClusterPool{}
		if runtime.DefaultUnstructuredConverter.FromUnstructured(cpu.UnstructuredContent(), cp); err != nil {
			return clusterPools, err
		}
		clusterPools.Items = append(clusterPools.Items, *cp)
	}
	return clusterPools, nil
}

func (cph *ClusterPoolHost) ConvertToPrintClusterPoolList(cpl *hivev1.ClusterPoolList, specificClusterPool string) (*printclusterpoolv1alpha1.PrintClusterPoolList, error) {
	pcps := &printclusterpoolv1alpha1.PrintClusterPoolList{}
	var singletonFound = false
	for i := range cpl.Items {
		// if only a specific cluster pool list is wanted, skip the others
		if specificClusterPool != "" {
			if specificClusterPool != cpl.Items[i].Name {
				continue
			} else {
				singletonFound = true
			}

		}
		pcp := printclusterpoolv1alpha1.PrintClusterPool{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cpl.Items[i].Name,
				Namespace: cpl.Items[i].Namespace,
			},
			Spec: printclusterpoolv1alpha1.PrintClusterPoolSpec{
				ClusterPoolHostName: cph.Name,
				ClusterPool:         &cpl.Items[i],
			},
		}
		pcps.Items = append(pcps.Items, pcp)
	}
	//If they only want one item, be sure we found that item or toss and error
	if specificClusterPool != "" && !singletonFound {
		return nil, fmt.Errorf("clusterpool %s was not found", specificClusterPool)
	}
	return pcps, nil
}

func (cph *ClusterPoolHost) GetClusterPoolConfig(clusterPoolName string, withoutCredentials bool, beta bool, outputFile string) error {
	values := make(map[string]interface{})
	reader := scenario.GetScenarioResourcesReader()

	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	apiExtensionsClient, err := apiextensionsclient.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return err
	}

	//Get clusterDeployment
	cpu, err := dynamicClient.Resource(helpers.GvrCP).Namespace(cph.Namespace).Get(context.TODO(), clusterPoolName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cp := &hivev1.ClusterPool{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(cpu.UnstructuredContent(), cp)
	if err != nil {
		return err
	}
	values["clusterPool"] = cpu.Object

	//Check if platform is supported.
	var secretName string
	switch {
	case cp.Spec.Platform.AWS != nil:
		secretName = cp.Spec.Platform.AWS.CredentialsSecretRef.Name
	case cp.Spec.Platform.Azure != nil && beta:
		secretName = cp.Spec.Platform.Azure.CredentialsSecretRef.Name
	case cp.Spec.Platform.GCP != nil && beta:
		secretName = cp.Spec.Platform.GCP.CredentialsSecretRef.Name
	default:
		return fmt.Errorf("unsupported platform %v", cp.Spec.Platform)
	}

	//Get install-config
	ic, err := kubeClient.CoreV1().Secrets(cph.Namespace).Get(context.TODO(), cp.Spec.InstallConfigSecretTemplateRef.Name, metav1.GetOptions{})
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

	//Get pull secret
	pk, err := kubeClient.CoreV1().Secrets(cph.Namespace).Get(context.TODO(), cp.Spec.PullSecretRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	values["imagePullSecret"] = string(pk.Data[".dockerconfigjson"])

	//Get credentials
	if withoutCredentials {
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
		cred, err := kubeClient.CoreV1().Secrets(cph.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		switch {
		case cp.Spec.Platform.AWS != nil:
			values["cloud"] = "aws"
			values["awsAccessKeyID"] = string(cred.Data["aws_access_key_id"])
			values["awsSecretAccessKey"] = string(cred.Data["aws_secret_access_key"])
		case cp.Spec.Platform.Azure != nil:
			values["cloud"] = "azure"
			osServicePrincipal := cred.Data["osServicePrincipal.json"]
			osServicePrincipalMap := make(map[string]interface{})
			err = json.Unmarshal(osServicePrincipal, &osServicePrincipalMap)
			if err != nil {
				return err
			}
			values["osServicePrincipalJson"] = osServicePrincipalMap
		case cp.Spec.Platform.GCP != nil:
			values["cloud"] = "gcp"
			values["osServiceAccountJson"] = string(cred.Data["osServiceAccount.json"])
		}
	}

	//Get clusterimageset
	klog.V(5).Infof("ImageSetRef:%s", cp.Spec.ImageSetRef.Name)
	values["imageSetRef"] = cp.Spec.ImageSetRef.Name

	klog.V(5).Infof("%v\n", values)
	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient)
	b, err := applier.MustTempalteAsset(reader, values, "", "config/clusterpool/config.yaml")
	if err != nil {
		return err
	}
	if len(outputFile) != 0 {
		return ioutil.WriteFile(outputFile, b, 0600)
	}
	fmt.Printf("%s\n", string(b))
	return nil
}
