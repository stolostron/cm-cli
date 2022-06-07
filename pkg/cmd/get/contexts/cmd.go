// Copyright Contributors to the Open Cluster Management project
package contexts

import (
	"fmt"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/helpers"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	example = `
	# get the contexts of the hub
	%[1]s get contexts
	
	# get the contexts and provide the clusterpoolhosts where the clusterclaim can be found
	%[1]s get contexts --cph <clusterpoolhosts>
`
)

// NewCmd ...
func NewCmd(f cmdutil.Factory, cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {

	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:                   "contexts",
		Aliases:               []string{"context"},
		DisableFlagsInUseLine: true,
		Short:                 "Get the managedcluster's contexts of a hub",
		Long:                  "Get the managedcluster's contexts of a hub based on hive clusterClaim and clusterDeployment",
		Example:               fmt.Sprintf(example, helpers.GetExampleHeader()),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run(streams))
		},
	}

	cmd.Flags().StringVar(&o.ClusterPoolHost, "cph", "", "The clusterpoolhost to use")

	return cmd
}
