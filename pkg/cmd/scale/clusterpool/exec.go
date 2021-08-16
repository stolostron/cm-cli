// Copyright Contributors to the Open Cluster Management project
package clusterpool

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("clusterclaim name is missing")
	}
	o.ClusterPool = args[0]
	return nil
}

func (o *Options) validate() error {
	return nil
}

func (o *Options) run() (err error) {
	cphs, err := clusterpoolhost.GetClusterPoolHosts()
	if err != nil {
		return err
	}

	cph, err := cphs.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	return cph.SizeClusterPool(o.ClusterPool, o.Size, o.CMFlags.DryRun)
}
