// Copyright Contributors to the Open Cluster Management project
package hypershiftdeployment

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
	# get the config for a hypershiftdeployment
	%[1]s get config hypershiftdeployment <cluster-name> -n <managedcluster_name> [--output-file <file_name]

`
)

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {

	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:                   "hypershiftdeployment",
		Aliases:               []string{"hypershiftdeployments", "hd", "hds"},
		DisableFlagsInUseLine: true,
		Short:                 "Display the config of a hypershiftdeployment",
		Example:               fmt.Sprintf(example, helpers.GetExampleHeader()),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run())
		},
	}

	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "The config will be copied in the specified file")

	return cmd
}
