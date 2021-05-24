// Copyright Contributors to the Open Cluster Management project
package hub

import (
	"github.com/open-cluster-management/cm-cli/pkg/cmd/applierscenarios"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	applierScenariosOptions *applierscenarios.ApplierScenariosOptions
	values                  map[string]interface{}
}

func newOptions(streams genericclioptions.IOStreams) *Options {
	return &Options{
		applierScenariosOptions: applierscenarios.NewApplierScenariosOptions(streams),
	}
}
