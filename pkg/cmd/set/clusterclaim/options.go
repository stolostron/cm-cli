// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	//CMFlags: The generic optiosn from the cm cli-runtime.
	CMFlags         *genericclioptionscm.CMFlags
	ClusterPoolHost string
	// list of clusterclaim names (comma-separated)
	ClusterClaims string
	// hibernate schedule on
	HibernateScheduleOn bool
	// hibernate schedule off
	HibernateScheduleOff bool
}

func newOptions(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		CMFlags: cmFlags,
	}
}
