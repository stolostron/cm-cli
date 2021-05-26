// Copyright Contributors to the Open Cluster Management project

package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/open-cluster-management/cm-cli/pkg/cmd/verbs"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	kubectlcmd "k8s.io/kubectl/pkg/cmd"
	cmdconfig "k8s.io/kubectl/pkg/cmd/config"
	"k8s.io/kubectl/pkg/cmd/options"
	"k8s.io/kubectl/pkg/cmd/plugin"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func main() {
	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	// root, f := clusteradmcmd.GetRootCmd("cm", streams)
	// root.AddCommand(newCmdCMVerbs(f, streams))

	// flags := pflag.NewFlagSet("cm", pflag.ExitOnError)
	// pflag.CommandLine = flags

	configFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(configFlags)
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	root := newCmdCMVerbs(f, streams)

	flags := root.PersistentFlags()
	matchVersionKubeConfigFlags.AddFlags(flags)
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc) // Warn for "_" flags

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	// From this point and forward we get warnings on flags that contain "_" separators
	root.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)

	configFlags.AddFlags(flags)
	root.AddCommand(cmdconfig.NewCmdConfig(f, clientcmd.NewDefaultPathOptions(), streams))
	//enable plugin functionality: all `os.Args[0]-<binary>` in the $PATH will be available for plugin
	plugin.ValidPluginFilenamePrefixes = []string{os.Args[0]}
	root.AddCommand(plugin.NewCmdPlugin(f, streams))
	root.AddCommand(kubectlcmd.NewDefaultKubectlCommand())
	root.AddCommand(options.NewCmdOptions(streams.Out))

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// NewCmdNamespace provides a cobra command wrapping NamespaceOptions
func newCmdCMVerbs(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{Use: "cm"}
	cmd.AddCommand(
		verbs.NewVerbCreate("create", f, streams),
		verbs.NewVerbGet("get", f, streams),
		verbs.NewVerbDelete("delete", f, streams),
		verbs.NewVerbApplier("applier", f, streams),
		verbs.NewVerbAttach("attach", f, streams),
		verbs.NewVerbDetach("detach", f, streams),
		verbs.NewVerbVersion("version", f, streams),
		verbs.NewVerbScale("scale", streams),
	)

	return cmd
}
