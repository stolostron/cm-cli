// Copyright Contributors to the Open Cluster Management project
package policies

import (
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/get"
)

type Options struct {
	//CMFlags: The generic options from the cm cli-runtime.
	CMFlags    *genericclioptionscm.CMFlags
	PolicyName string
	GetOptions *get.GetOptions
}

func newOptions(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		CMFlags:    cmFlags,
		GetOptions: get.NewGetOptions("cm", streams),
	}
}
