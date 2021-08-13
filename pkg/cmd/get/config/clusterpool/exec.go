// Copyright Contributors to the Open Cluster Management project
package clusterpool

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 0 {
		o.ClusterPoolName = args[0]
	}
	return nil
}

func (o *Options) validate() error {
	if len(o.ClusterPoolName) == 0 {
		return fmt.Errorf("clusterpoolname is missing")
	}
	return nil
}

func (o *Options) run() (err error) {

	cphs, err := clusterpoolhost.GetClusterPoolHosts()
	if err != nil {
		return err
	}

	cph, err := cphs.GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}

	if len(o.ClusterPoolHost) != 0 {
		cph, err = cphs.GetClusterPoolHost(o.ClusterPoolHost)
		if err != nil {
			return err
		}
	}

	return cph.GetClusterPoolConfig(o.ClusterPoolName, o.withoutCredentials, o.CMFlags.Beta, o.outputFile)

}
