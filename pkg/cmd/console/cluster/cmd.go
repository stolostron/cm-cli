// Copyright Contributors to the Open Cluster Management project
package cluster

import (
	"fmt"

	genericclioptionscm "github.com/stolostron/cm-cli/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/helpers"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	example = `
	# open the console of a managedcluster cluster
	%[1]s concole cluster <managed-cluster-name> 
`
)

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {

	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:                   "cluster",
		Aliases:               []string{"clusters"},
		DisableFlagsInUseLine: true,
		Short:                 "Open console of a managed cluster",
		Long:                  "Open the console of a managed cluster if the managed cluster was created using a hypershift or cluster claim",
		Example:               fmt.Sprintf(example, helpers.GetExampleHeader()),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run())
		},
	}

	o.GetOptions.PrintFlags = get.NewGetPrintFlags()

	o.GetOptions.PrintFlags.AddFlags(cmd)

	cmd.Flags().StringVar(&o.ClusterPoolHost, "cph", "", "The clusterpoolhost to use")
	cmd.Flags().IntVar(&o.Timeout, "timeout", 60, "Timeout to get the cluster claim running")
	cmd.Flags().BoolVar(&o.WithCredentials, "creds", o.WithCredentials, "If set the credentials will be displayed")

	return cmd
}
