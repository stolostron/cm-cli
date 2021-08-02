// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/get"
)

type Options struct {
	//CMFlags: The generic optiosn from the cm cli-runtime.
	CMFlags             *genericclioptionscm.CMFlags
	ClusterClaim        string
	PrintFlags          *get.PrintFlags
	AllClusterPoolHosts bool
	ClusterPoolHost     string
	Timeout             int
}

func newOptions(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		CMFlags: cmFlags,
	}
}
