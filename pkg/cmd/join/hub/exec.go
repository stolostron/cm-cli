// Copyright Contributors to the Open Cluster Management project
package hub

import (
	"fmt"
	"path/filepath"

	"github.com/ghodss/yaml"
	appliercmd "github.com/open-cluster-management/applier/pkg/applier/cmd"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/join/hub/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api/v1"

	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		return nil
	}
	//Check if default values must be used
	if o.applierScenariosOptions.ValuesPath == "" {
		if o.token != "" && o.hubServer != "" {
			o.values = make(map[string]interface{})
			hub := make(map[string]interface{})
			hub["token"] = o.token
			hub["hubServer"] = o.hubServer
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

	if o.hubServer == "" {
		ihubServer, ok := hub["hubServer"]
		if !ok || ihubServer == nil {
			return fmt.Errorf("hub-server name is missing")
		}
		o.token = ihubServer.(string)
	}

	hub["hubServer"] = o.hubServer

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

	bootstrapConfig := clientcmdapi.Config{
		// Define a cluster stanza based on the bootstrap kubeconfig.
		Clusters: []clientcmdapi.NamedCluster{
			{
				Name: "hub",
				Cluster: clientcmdapi.Cluster{
					Server:                o.hubServer,
					InsecureSkipTLSVerify: true,
				},
			},
		},
		// Define auth based on the obtained client cert.
		AuthInfos: []clientcmdapi.NamedAuthInfo{
			{
				Name: "bootstrap",
				AuthInfo: clientcmdapi.AuthInfo{
					Token: string(o.token),
				},
			},
		},
		// Define a context that connects the auth info and cluster, and set it as the default
		Contexts: []clientcmdapi.NamedContext{
			{
				Name: "bootstrap",
				Context: clientcmdapi.Context{
					Cluster:   "hub",
					AuthInfo:  "bootstrap",
					Namespace: "default",
				},
			},
		},
		CurrentContext: "bootstrap",
	}
	bootstrapConfigBytes, err := yaml.Marshal(bootstrapConfig)
	if err != nil {
		return err
	}

	o.values["kubeconfig"] = string(bootstrapConfigBytes)
	fmt.Printf("%v", o.values)

	err = applyOptions.ApplyWithValues(client, reader,
		filepath.Join(scenarioDirectory, "hub"), []string{},
		o.values)

	if err != nil {
		return err
	}

	fmt.Printf("login back onto the hub and run:  ")

	return nil

}
