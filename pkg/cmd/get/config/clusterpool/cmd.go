// Copyright Contributors to the Open Cluster Management project
package clusterpool

import (
	"fmt"

	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	example = `
	# get the config for a cluster
	%[1]s get config clusterpool <clusterpool_name> [--output-file <file_name]

`
)

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {

	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:                   "clusterpool",
		Aliases:               []string{"clusterpools", "cp", "cphs"},
		DisableFlagsInUseLine: true,
		Short:                 "Display the config of a cluster",
		Example:               fmt.Sprintf(example, helpers.GetExampleHeader()),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run())
		},
	}

	cmd.Flags().StringVar(&o.ClusterPoolHost, "cph", "", "The clusterpoolhost to use")
	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "The config will be copied in the specified file")
	cmd.Flags().BoolVar(&o.withoutCredentials, "without-credentials", false, "If set the platform credentials will be not inserted in the config")

	return cmd
}
