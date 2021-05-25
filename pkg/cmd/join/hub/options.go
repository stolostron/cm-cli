// Copyright Contributors to the Open Cluster Management project
package hub

import (
	"github.com/open-cluster-management/cm-cli/pkg/cmd/applierscenarios"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type Options struct {
	applierScenariosOptions *applierscenarios.ApplierScenariosOptions
	token                   string
	hubServerInternal       string
	hubServerExternal       string
	clusterName             string
	factory                 cmdutil.Factory
	values                  map[string]interface{}
}

func newOptions(f cmdutil.Factory, streams genericclioptions.IOStreams) *Options {
	return &Options{
		applierScenariosOptions: applierscenarios.NewApplierScenariosOptions(streams),
		factory:                 f,
	}
}
