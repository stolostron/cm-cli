// Copyright Contributors to the Open Cluster Management project
package components

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
	# get the available components
	%[1]s get components
`
)

// NewCmd ...
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) *cobra.Command {

	o := newOptions(cmFlags, streams)
	cmd := &cobra.Command{
		Use:                   "components",
		Aliases:               []string{"component"},
		DisableFlagsInUseLine: true,
		Short:                 "Get the list of available components",
		Example:               fmt.Sprintf(example, helpers.GetExampleHeader()),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.complete(cmd, args))
			cmdutil.CheckErr(o.validate())
			cmdutil.CheckErr(o.run(streams))
		},
	}

	return cmd
}
