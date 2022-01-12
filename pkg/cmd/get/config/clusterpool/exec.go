// Copyright Contributors to the Open Cluster Management project
package clusterpool

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"
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

	cph, err := clusterpoolhost.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	return cph.GetClusterPoolConfig(o.ClusterPoolName, o.withoutCredentials, o.CMFlags.Beta, o.outputFile)

}
