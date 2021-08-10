// Copyright Contributors to the Open Cluster Management project

package main

import (
	"flag"
	"os"

	"github.com/spf13/cobra"

	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	cmdconfig "k8s.io/kubectl/pkg/cmd/config"
	"k8s.io/kubectl/pkg/cmd/options"
	"k8s.io/kubectl/pkg/cmd/plugin"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"

	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"

	"github.com/open-cluster-management/cm-cli/pkg/cmd/accept"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/attach"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/console"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/create"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/delete"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/detach"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/enable"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/get"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/hibernate"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/initialization"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/join"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/run"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/scale"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/set"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/use"
	"github.com/open-cluster-management/cm-cli/pkg/cmd/version"
)

func main() {
	root := &cobra.Command{
		Use: "cm",
	}

	flags := root.PersistentFlags()
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc) // Warn for "_" flags

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	kubeConfigFlags.WrapConfigFn = setQPS
	kubeConfigFlags.AddFlags(flags)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)

	matchVersionKubeConfigFlags.AddFlags(flags)

	klog.InitFlags(nil)
	flags.AddGoFlagSet(flag.CommandLine)

	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	// From this point and forward we get warnings on flags that contain "_" separators
	root.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)
	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	clusteradmFlags := genericclioptionsclusteradm.NewClusteradmFlags(f)

	cmFlags := genericclioptionscm.NewCMFlags(f)
	cmFlags.AddFlags(flags)

	root.AddCommand(cmdconfig.NewCmdConfig(f, clientcmd.NewDefaultPathOptions(), streams))
	root.AddCommand(options.NewCmdOptions(streams.Out))

	//enable plugin functionality: all `os.Args[0]-<binary>` in the $PATH will be available for plugin
	plugin.ValidPluginFilenamePrefixes = []string{os.Args[0]}
	root.AddCommand(plugin.NewCmdPlugin(f, streams))

	groups := templates.CommandGroups{
		{
			Message: "General commands:",
			Commands: []*cobra.Command{
				version.NewCmd(cmFlags, streams),
			},
		},
		{
			Message: "Clusters commands:",
			Commands: []*cobra.Command{
				attach.NewCmd(cmFlags, streams),
				detach.NewCmd(cmFlags, streams),
				create.NewCmd(cmFlags, streams),
				delete.NewCmd(clusteradmFlags, cmFlags, streams),
				scale.NewCmd(cmFlags, streams),
				enable.NewCmd(cmFlags, streams),
				get.NewCmd(f, clusteradmFlags, cmFlags, streams),
				initialization.NewCmd(clusteradmFlags, cmFlags, streams),
			},
		},
		{
			Message: "Registration commands:",
			Commands: []*cobra.Command{
				join.NewCmd(clusteradmFlags, cmFlags, streams),
				accept.NewCmd(clusteradmFlags, cmFlags, streams),
			},
		},
		{
			Message: "Cluster pools commands:",
			Commands: []*cobra.Command{
				use.NewCmd(cmFlags, streams),
				set.NewCmd(cmFlags, streams),
				run.NewCmd(cmFlags, streams),
				hibernate.NewCmd(cmFlags, streams),
				console.NewCmd(cmFlags, streams),
			},
		},
		{
			Message: "Governance/policy commands:",
			Commands: []*cobra.Command{
				get.NewCmd(f, clusteradmFlags, cmFlags, streams),
			},
		},
	}
	groups.Add(root)
	err := root.Execute()
	if err != nil {
		klog.V(1).ErrorS(err, "Error:")
	}
	klog.Flush()
	if err != nil {
		os.Exit(1)
	}
}

func setQPS(r *rest.Config) *rest.Config {
	r.QPS = helpers.QPS
	r.Burst = helpers.Burst
	return r
}
