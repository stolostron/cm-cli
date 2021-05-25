// Copyright Contributors to the Open Cluster Management project
package clusters

import (
	"github.com/open-cluster-management/cm-cli/pkg/cmd/applierscenarios"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	applierScenariosOptions *applierscenarios.ApplierScenariosOptions
	factory                 cmdutil.Factory
	clusters                string
	values                  map[string]interface{}
}

func newOptions(f cmdutil.Factory, streams genericclioptions.IOStreams) *Options {
	return &Options{
		applierScenariosOptions: applierscenarios.NewApplierScenariosOptions(streams),
		factory:                 f,
	}
}
