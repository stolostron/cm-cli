// Copyright Contributors to the Open Cluster Management project
package hub

import (
	"fmt"
	"path/filepath"

	"github.com/ghodss/yaml"
	appliercmd "github.com/open-cluster-management/applier/pkg/applier/cmd"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/join/hub/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"

	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		return nil
	}
	//Check if default values must be used
	if o.applierScenariosOptions.ValuesPath == "" {
		if o.token != "" && o.hubServerInternal != "" {
			o.values = make(map[string]interface{})
			hub := make(map[string]interface{})
			hub["token"] = o.token
			hub["hubServerInternal"] = o.hubServerInternal
			hub["hubServerExternal"] = o.hubServerExternal
			hub["clusterName"] = o.clusterName
			o.values["hub"] = hub
		} else {
			return fmt.Errorf("values or token/hub-server are missing")
		}
	} else {
		//Read values
		o.values, err = appliercmd.ConvertValuesFileToValuesMap(o.applierScenariosOptions.ValuesPath, "")
		if err != nil {
			return err
		}
	}

	if len(o.values) == 0 {
		return fmt.Errorf("values are missing")
	}

	return nil
}

func (o *Options) validate() error {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		return nil
	}
	ihub, ok := o.values["hub"]
	if !ok || ihub == nil {
		return fmt.Errorf("hub is missing")
	}
	hub := ihub.(map[string]interface{})

	if o.token == "" {
		itoken, ok := hub["token"]
		if !ok || itoken == nil {
			return fmt.Errorf("token name is missing")
		}
		o.token = itoken.(string)
	}

	hub["token"] = o.token

	if o.hubServerInternal == "" {
		ihubServer, ok := hub["hubServerInternal"]
		if !ok || ihubServer == nil {
			return fmt.Errorf("hub-server-internal name is missing")
		}
		o.token = ihubServer.(string)
	}

	hub["hubServerInternal"] = o.hubServerInternal

	if o.hubServerExternal == "" {
		ihubServer, ok := hub["hubServerExternal"]
		if !ok || ihubServer == nil {
			return fmt.Errorf("hub-server-external name is missing")
		}
		o.token = ihubServer.(string)
	}

	hub["hubServerExternal"] = o.hubServerExternal

	if o.clusterName == "" {
		iclusterName, ok := hub["clusterName"]
		if !ok || iclusterName == nil {
			return fmt.Errorf("name is missing")
		}
		o.token = iclusterName.(string)
	}

	hub["clusterName"] = o.clusterName

	return nil
}

func (o *Options) run() error {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		reader := scenario.GetApplierScenarioResourcesReader()
		return reader.ExtractAssets(scenarioDirectory, o.applierScenariosOptions.OutTemplatesDir)
	}
	client, err := helpers.GetControllerRuntimeClientFromFlags(o.applierScenariosOptions.ConfigFlags)
	if err != nil {
		return err
	}
	return o.runWithClient(client)
}

func (o *Options) runWithClient(client crclient.Client) error {
	reader := scenario.GetApplierScenarioResourcesReader()

	applyOptions := &appliercmd.Options{
		OutFile:     o.applierScenariosOptions.OutFile,
		ConfigFlags: o.applierScenariosOptions.ConfigFlags,

		Delete:    false,
		Timeout:   o.applierScenariosOptions.Timeout,
		Force:     o.applierScenariosOptions.Force,
		Silent:    o.applierScenariosOptions.Silent,
		IOStreams: o.applierScenariosOptions.IOStreams,
	}

	bootstrapConfigUnSecure := clientcmdapiv1.Config{
		// Define a cluster stanza based on the bootstrap kubeconfig.
		Clusters: []clientcmdapiv1.NamedCluster{
			{
				Name: "hub",
				Cluster: clientcmdapiv1.Cluster{
					Server:                o.hubServerExternal,
					InsecureSkipTLSVerify: true,
				},
			},
		},
		// Define auth based on the obtained client cert.
		AuthInfos: []clientcmdapiv1.NamedAuthInfo{
			{
				Name: "bootstrap",
				AuthInfo: clientcmdapiv1.AuthInfo{
					Token: string(o.token),
				},
			},
		},
		// Define a context that connects the auth info and cluster, and set it as the default
		Contexts: []clientcmdapiv1.NamedContext{
			{
				Name: "bootstrap",
				Context: clientcmdapiv1.Context{
					Cluster:   "hub",
					AuthInfo:  "bootstrap",
					Namespace: "default",
				},
			},
		},
		CurrentContext: "bootstrap",
	}

	bootstrapConfigBytesUnSecure, err := yaml.Marshal(bootstrapConfigUnSecure)
	if err != nil {
		return err
	}

	configUnSecure, err := clientcmd.Load(bootstrapConfigBytesUnSecure)
	if err != nil {
		return err
	}
	restConfigUnSecure, err := clientcmd.NewDefaultClientConfig(*configUnSecure, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return err
	}

	clientUnSecure, err := crclient.New(restConfigUnSecure, crclient.Options{})
	if err != nil {
		return err
	}

	ca, err := helpers.GetCACert(clientUnSecure)
	if err != nil {
		return err
	}

	bootstrapConfig := bootstrapConfigUnSecure
	bootstrapConfig.Clusters[0].Cluster.InsecureSkipTLSVerify = false
	bootstrapConfig.Clusters[0].Cluster.CertificateAuthorityData = ca
	bootstrapConfig.Clusters[0].Cluster.Server = o.hubServerInternal
	bootstrapConfigBytes, err := yaml.Marshal(bootstrapConfig)
	if err != nil {
		return err
	}

	o.values["kubeconfig"] = string(bootstrapConfigBytes)

	restConfig, err := o.factory.ToRESTConfig()
	if err != nil {
		return err
	}

	o.values["clusterServerExternal"] = restConfig.Host

	err = applyOptions.ApplyWithValues(client, reader,
		filepath.Join(scenarioDirectory, "hub"), []string{},
		o.values)

	if err != nil {
		return err
	}

	fmt.Printf("login back onto the hub and run: %s accept clusters --names %s\n", helpers.GetExampleHeader(), o.clusterName)

	return nil

}
