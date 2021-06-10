// Copyright Contributors to the Open Cluster Management project
package version

import (
	"fmt"

	genericclioptionscm "github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Version
%[1]s version
`

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(cmFlags *genericclioptionscm.CMFlags, streams genericclioptions.IOStreams) (cmd *cobra.Command) {
	o := newOptions(cmFlags, streams)
	cmd = &cobra.Command{
		Use:          "version",
		Short:        "get the versions of the different components",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PreRunE: func(c *cobra.Command, args []string) error {
			if !helpers.IsRHACM(cmFlags.KubectlFactory) {
				return fmt.Errorf("this command '%s version' is only available on RHACM", helpers.GetExampleHeader())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.complete(c, args); err != nil {
				return err
			}
			if err := o.validate(); err != nil {
				return err
			}
			if err := o.run(); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}
