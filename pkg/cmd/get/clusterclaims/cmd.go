// Copyright Contributors to the Open Cluster Management project
package clusterclaims

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	example = `
	# get a clusterclaim in current clusterpoolhost
	%[1]s clusterclaims|cc <cluster-name> <clusterpoolhosts>`
)

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {

	o := newOptions(cmFlags, streams)
	clusters := &cobra.Command{
		Use:                   "clusterclaims",
		Aliases:               []string{"clusterclaim", "cc"},
		DisableFlagsInUseLine: true,
		Short:                 "Display the clusterclaims",
		Example:               fmt.Sprintf(example, helpers.GetExampleHeader()),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return clusterpoolhost.BackupCurrentContexts()
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run())
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return clusterpoolhost.RestoreCurrentContexts()
		},
		SuggestFor: []string{"list", "ps"},
	}

	clusters.Flags().BoolVarP(&o.AllClusterPoolHosts, "all-cphs", "A", o.AllClusterPoolHosts, "If the requested object does not exist the command will return exit code 0.")
	return clusters
}
