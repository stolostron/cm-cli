// Copyright Contributors to the Open Cluster Management project
package genericclioptions

import (
	"github.com/spf13/pflag"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type CMFlags struct {
	KubectlFactory cmdutil.Factory
	//if set the resources will be sent to stdout instead of being applied
	DryRun bool
	//Accept beta cmd
	Beta bool
	//Namespace where the server (RHACM/MCE) is installed)
	ServerNamespace string
	//if set cm will not check if the server (RHACM/MCE) is installed
	SkipServerCheck bool
}

// NewClusteradmFlags returns CMFlags with default values set
func NewCMFlags(f cmdutil.Factory) *CMFlags {
	return &CMFlags{
		KubectlFactory: f,
	}
}

func (f *CMFlags) AddFlags(flags *pflag.FlagSet) {
	flags.BoolVar(&f.DryRun, "dry-run", false, "If set the generated resources will be displayed but not applied")
	flags.BoolVar(&f.Beta, "beta", false, "If set commands or functionalities in beta version will be available")
	flags.StringVar(&f.ServerNamespace, "server-namespace", "", "The namespace where the server (RHACM/MCE) is installed")
	flags.BoolVar(&f.SkipServerCheck, "skip-server-check", false, "If set commands will not check the installed server (RHACM/MCE) target")
}
