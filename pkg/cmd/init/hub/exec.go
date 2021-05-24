// Copyright Contributors to the Open Cluster Management project
package hub

import (
	"fmt"
	"path/filepath"

	appliercmd "github.com/open-cluster-management/applier/pkg/applier/cmd"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/init/hub/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"

	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		return nil
	}
	//Check if default values must be used
	if o.applierScenariosOptions.ValuesPath == "" {
		o.values = make(map[string]interface{})
		hub := make(map[string]interface{})
		hub["tokenID"] = helpers.RandStringRunes_az09(6)
		hub["tokenSecret"] = helpers.RandStringRunes_az09(16)
		o.values["hub"] = hub
	} else {
		//Read values
		o.values, err = appliercmd.ConvertValuesFileToValuesMap(o.applierScenariosOptions.ValuesPath, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Options) validate() error {
	if o.applierScenariosOptions.OutTemplatesDir != "" {
		return nil
	}
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

	err := applyOptions.ApplyWithValues(client, reader,
		filepath.Join(scenarioDirectory, "hub"), []string{},
		o.values)

	if err != nil {
		return err
	}

	apiServer, err := helpers.GetAPIServer(client)
	if err != nil {
		return err
	}

	fmt.Printf("login into the cluster and run: %s join hub --hub-token %s.%s --hub-server %s\n",
		helpers.GetExampleHeader(),
		o.values["hub"].(map[string]interface{})["tokenID"].(string),
		o.values["hub"].(map[string]interface{})["tokenSecret"].(string),
		apiServer,
	)

	return nil
}
