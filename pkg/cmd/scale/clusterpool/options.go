// Copyright Contributors to the Open Cluster Management project
package clusterpool

import (
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	//CMFlags: The generic optiosn from the cm cli-runtime.
	CMFlags *genericclioptionscm.CMFlags
	//The list of cluster claim name to create (comma-separated)
	ClusterPool     string
	ClusterPoolHost string
	Size            int32
}

func newOptions(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		CMFlags: cmFlags,
	}
}
