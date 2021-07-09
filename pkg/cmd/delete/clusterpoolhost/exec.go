// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("clusterpoolcph name is missing")
	}
	o.ClusterHostPool = args[0]
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
	cph, err := cphs.GetClusterPoolHost(o.ClusterHostPool)
	if err != nil {
		return err
	}
	err = cph.DeleteClusterPoolHost()
	if err != nil {
		return err
	}
	return nil
}
