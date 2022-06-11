// Copyright Contributors to the Open Cluster Management project
package components

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/helpers"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 0 {
		o.component = args[0]
	}
	if len(o.component) == 0 {
		return fmt.Errorf("component name is missing")
	}
	return nil
}

func (o *Options) validate() error {
	return nil
}

func (o *Options) run(streams genericclioptions.IOStreams) (err error) {
	return helpers.SetComponentEnable(o.CMFlags, o.component, true)
}
