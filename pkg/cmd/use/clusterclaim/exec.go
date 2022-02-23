// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("clustername is missing")
	}
	o.Cluster = args[0]
	return nil
}

func (o *Options) validate() error {
	return nil
}

func (o *Options) run() (err error) {
	cph, err := clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	return cph.SetClusterClaimContext(o.Cluster, true, o.Timeout, o.CMFlags.DryRun, o.outputFile, nil)
}
