// Copyright Contributors to the Open Cluster Management project
package hypershiftdeployment

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/stolostron/applier/pkg/apply"
	"github.com/stolostron/cm-cli/pkg/cmd/get/config/hypershiftdeployment/scenario"
	"github.com/stolostron/cm-cli/pkg/hypershift"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 0 {
		o.clusterName = args[0]
	}
	return nil
}

func (o *Options) validate() error {
	if len(o.clusterName) == 0 {
		return fmt.Errorf("hypershiftdeployment name is missing")
	}
	return nil
}

func (o *Options) run() (err error) {
	restConfig, err := o.CMFlags.KubectlFactory.ToRESTConfig()
	if err != nil {
		return err
	}
	return o.runWithClient(restConfig)
}

func (o *Options) runWithClient(restConfig *rest.Config) (err error) {
	reader := scenario.GetScenarioResourcesReader()

	//Get hypershiftdeployment
	hd, err := hypershift.GetHypershiftDeployment(o.clusterName, o.CMFlags)
	if err != nil {
		return err
	}

	klog.V(5).Infof("%v\n", hd)
	applierBuilder := apply.NewApplierBuilder()
	applier := applierBuilder.WithRestConfig(restConfig).Build()
	b, err := applier.MustTemplateAsset(reader, *hd, "", "config/config.yaml")
	if err != nil {
		return err
	}
	if len(o.outputFile) != 0 {
		return ioutil.WriteFile(o.outputFile, b, 0600)
	}
	fmt.Printf("%s\n", string(b))
	return nil
}
