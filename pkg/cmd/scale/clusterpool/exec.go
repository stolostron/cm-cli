// Copyright Contributors to the Open Cluster Management project
package clusterpool

import (
	"fmt"

	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("clusterpool name is missing")
	}
	o.ClusterPool = args[0]

	if !cmd.Flags().Changed("size") {
		return fmt.Errorf("size must be specified")
	}
	return nil
}

func (o *Options) validate() error {
	if o.Size < 0 {
		return fmt.Errorf("size must be greater than or equal to zero")
	}
	return nil
}

func (o *Options) run() (err error) {
	cph, err := clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	return cph.SizeClusterPool(o.ClusterPool, o.Size, o.CMFlags.DryRun)
}
