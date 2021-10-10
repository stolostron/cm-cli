// Copyright Contributors to the Open Cluster Management project
package authrealm

import (
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	//CMFlags: The generic optiosn from the cm cli-runtime.
	CMFlags               *genericclioptionscm.CMFlags
	name                  string
	namespace             string
	typeName              string
	routeSubDomain        string
	placementName         string
	clusterSetName        string
	clusterSetBindingName string
	valuesPath            string
	values                map[string]interface{}
	//The file to output the resources will be sent to the file.
	outputFile string
}

func newOptions(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		CMFlags: cmFlags,
	}
}
