// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	//CMFlags: The generic optiosn from the cm cli-runtime.
	CMFlags                  *genericclioptionscm.CMFlags
	valuesPath               string
	values                   map[string]interface{}
	clusterName              string
	clusterServer            string
	clusterToken             string
	clusterKubeConfig        string
	clusterKubeConfigContent string
	importFile               string
	waitAgent                bool
	waitAddOns               bool
	timeout                  int
	hiveScenario             bool

	//The file to output the resources will be sent to the file.
	outputFile string
}

func newOptions(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		CMFlags: cmFlags,
	}
}
