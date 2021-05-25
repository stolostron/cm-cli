// Copyright Contributors to the Open Cluster Management project
package hub

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	appliercmd "github.com/open-cluster-management/applier/pkg/applier/cmd"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/init/hub/scenario"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"k8s.io/apimachinery/pkg/labels"

	corev1 "k8s.io/api/core/v1"
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
	ss := &corev1.SecretList{}
	ls := labels.SelectorFromSet(labels.Set{
		"app": "cluster-manager",
	})
	err := client.List(context.TODO(),
		ss,
		&crclient.ListOptions{
			LabelSelector: ls,
			Namespace:     "kube-system",
		})
	if err != nil {
		return err
	}
	var bootstrapSecret *corev1.Secret
	for _, item := range ss.Items {
		if strings.HasPrefix(item.Name, "bootstrap-token") {
			bootstrapSecret = &item
			break
		}
	}
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

	if bootstrapSecret == nil {
		err = applyOptions.ApplyWithValues(client, reader,
			filepath.Join(scenarioDirectory, "hub"), []string{},
			o.values)
	} else {
		o.values["hub"].(map[string]interface{})["tokenID"] = string(bootstrapSecret.Data["token-id"])
		o.values["hub"].(map[string]interface{})["tokenSecret"] = string(bootstrapSecret.Data["token-secret"])
		err = applyOptions.ApplyWithValues(client, reader,
			filepath.Join(scenarioDirectory, "hub"), []string{"boostrap-token-secret.yaml"},
			o.values)
	}

	if err != nil {
		return err
	}

	apiServerInternal, err := helpers.GetAPIServer(client)
	if err != nil {
		return err
	}

	restConfig, err := o.factory.ToRESTConfig()
	if err != nil {
		return err
	}

	fmt.Printf("login into the cluster and run: %s join hub --hub-token %s.%s --hub-server-internal %s --hub-server-external %s --name <cluster_name>\n",
		helpers.GetExampleHeader(),
		o.values["hub"].(map[string]interface{})["tokenID"].(string),
		o.values["hub"].(map[string]interface{})["tokenSecret"].(string),
		apiServerInternal,
		restConfig.Host,
	)

	return nil
}
