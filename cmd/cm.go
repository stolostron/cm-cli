// Copyright Contributors to the Open Cluster Management project

package main

import (
	"flag"
	"os"

	"github.com/spf13/cobra"

	"github.com/open-cluster-management/cm-cli/pkg/cmd/version"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	cmdconfig "k8s.io/kubectl/pkg/cmd/config"
	"k8s.io/kubectl/pkg/cmd/options"
	"k8s.io/kubectl/pkg/cmd/plugin"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"

	// genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"

	// clusteradmaccept "open-cluster-management.io/clusteradm/pkg/cmd/accept"
	// clusteradminit "open-cluster-management.io/clusteradm/pkg/cmd/init"
	// clusteradmjoin "open-cluster-management.io/clusteradm/pkg/cmd/join"

	attachcluster "github.com/open-cluster-management/cm-cli/pkg/cmd/attach/cluster"
	createcluster "github.com/open-cluster-management/cm-cli/pkg/cmd/create/cluster"
	deletecluster "github.com/open-cluster-management/cm-cli/pkg/cmd/delete/cluster"
	detachcluster "github.com/open-cluster-management/cm-cli/pkg/cmd/detach/cluster"
	getclusters "github.com/open-cluster-management/cm-cli/pkg/cmd/get/clusters"
)

func main() {
	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	configFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(configFlags)
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	// clusteradmFlags := genericclioptionsclusteradm.NewClusteradmFlags(f)
	cmFlags := genericclioptionscm.NewCMFlags(f)
	root := &cobra.Command{
		Use: "cm",
	}
	// root := newCmdCMVerbs(f, streams)

	flags := root.PersistentFlags()
	matchVersionKubeConfigFlags.AddFlags(flags)
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc) // Warn for "_" flags

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	// From this point and forward we get warnings on flags that contain "_" separators
	root.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)

	configFlags.AddFlags(flags)
	cmFlags.AddFlags(flags)
	flags.AddGoFlagSet(flag.CommandLine)
	root.AddCommand(cmdconfig.NewCmdConfig(f, clientcmd.NewDefaultPathOptions(), streams))
	//enable plugin functionality: all `os.Args[0]-<binary>` in the $PATH will be available for plugin
	plugin.ValidPluginFilenamePrefixes = []string{os.Args[0]}
	root.AddCommand(plugin.NewCmdPlugin(f, streams))
	root.AddCommand(options.NewCmdOptions(streams.Out))
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
				attachcluster.NewCmd(cmFlags, streams),
				detachcluster.NewCmd(cmFlags, streams),
				createcluster.NewCmd(cmFlags, streams),
				deletecluster.NewCmd(cmFlags, streams),
				getclusters.NewCmd(cmFlags, streams),
			},
		},
		// {
		// 	Message: "Registration commands:",
		// 	Commands: []*cobra.Command{
		// 		clusteradminit.NewCmd(clusteradmFlags, streams),
		// 		clusteradmjoin.NewCmd(clusteradmFlags, streams),
		// 		clusteradmaccept.NewCmd(clusteradmFlags, streams),
		// 	},
		// },
	}
	groups.Add(root)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
