// Copyright Contributors to the Open Cluster Management project

package cmd

import (
	"flag"
	"os"

	"github.com/spf13/cobra"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"github.com/stolostron/cm-cli/pkg/helpers"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/cmd/options"
	"k8s.io/kubectl/pkg/cmd/plugin"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"

	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"

	"github.com/stolostron/cm-cli/pkg/cmd/accept"
	"github.com/stolostron/cm-cli/pkg/cmd/attach"
	"github.com/stolostron/cm-cli/pkg/cmd/console"
	"github.com/stolostron/cm-cli/pkg/cmd/create"
	"github.com/stolostron/cm-cli/pkg/cmd/delete"
	"github.com/stolostron/cm-cli/pkg/cmd/detach"
	"github.com/stolostron/cm-cli/pkg/cmd/enable"
	"github.com/stolostron/cm-cli/pkg/cmd/get"
	"github.com/stolostron/cm-cli/pkg/cmd/hibernate"
	"github.com/stolostron/cm-cli/pkg/cmd/initialization"
	"github.com/stolostron/cm-cli/pkg/cmd/install"
	"github.com/stolostron/cm-cli/pkg/cmd/join"
	"github.com/stolostron/cm-cli/pkg/cmd/run"
	"github.com/stolostron/cm-cli/pkg/cmd/scale"
	"github.com/stolostron/cm-cli/pkg/cmd/set"
	"github.com/stolostron/cm-cli/pkg/cmd/use"
	"github.com/stolostron/cm-cli/pkg/cmd/version"
	"github.com/stolostron/cm-cli/pkg/cmd/with"
)

func NewCMCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "cm",
		Short: "CLI for Red Hat Advanced Cluster Management",
		//This remove the auto-generated tag in the cobra doc
		DisableAutoGenTag: true,
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

	// root.AddCommand(cmdconfig.NewCmdConfig(f, clientcmd.NewDefaultPathOptions(), streams))
	root.AddCommand(options.NewCmdOptions(streams.Out))

	//enable plugin functionality: all `os.Args[0]-<binary>` in the $PATH will be available for plugin
	plugin.ValidPluginFilenamePrefixes = []string{os.Args[0]}
	root.AddCommand(plugin.NewCmdPlugin(streams))

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
			Message: "cluster pools commands:",
			Commands: []*cobra.Command{
				use.NewCmd(cmFlags, streams),
				set.NewCmd(cmFlags, streams),
				run.NewCmd(cmFlags, streams),
				hibernate.NewCmd(cmFlags, streams),
				console.NewCmd(cmFlags, streams),
				with.NewCmd(cmFlags, streams),
			},
		},
		{
			Message: "install commands:",
			Commands: []*cobra.Command{
				install.NewCmd(cmFlags, streams),
			},
		},
	}
	groups.Add(root)
	return root
}

func setQPS(r *rest.Config) *rest.Config {
	r.QPS = helpers.QPS
	r.Burst = helpers.Burst
	return r
}
