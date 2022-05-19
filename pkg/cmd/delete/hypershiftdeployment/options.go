// Copyright Contributors to the Open Cluster Management project
package hypershiftdeployment

import (
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	//CMFlags: The generic optiosn from the cm cli-runtime.
	CMFlags *genericclioptionscm.CMFlags
	//The list of hypershiftdeployment to delete (comma-separated)
	HypershiftDeployments string
	//The namespace of all hypershiftdeployemnt
	HypershiftDeploymentsNamespace string
}

func newOptions(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		CMFlags: cmFlags,
	}
}
