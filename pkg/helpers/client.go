// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"context"
	"fmt"

	"github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func GetControllerRuntimeClientFromFlags(configFlags *genericclioptions.ConfigFlags) (client crclient.Client, err error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	config.QPS = 20
	return crclient.New(config, crclient.Options{})
}

func GetAPIServer(client crclient.Client) (string, error) {
	cm := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), crclient.ObjectKey{Namespace: "kube-public", Name: "cluster-info"}, cm)
	if err != nil {
		return "", err
	}
	config := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(cm.Data["kubeconfig"]), &config)
	if err != nil {
		return "", err
	}
	clusters, ok := config["clusters"].([]interface{})
	if !ok || len(clusters) != 1 {
		return "", fmt.Errorf("can not find the cluster in the cluster-info")
	}
	cluster0, ok := clusters[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("can not find the cluster")
	}
	cluster, ok := cluster0["cluster"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("cluster not found")
	}
	return cluster["server"].(string), nil
}
