// Copyright Contributors to the Open Cluster Management project
package clusterpoolhosts

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(o.OutputFormat) == 0 {
		o.OutputFormat = helpers.CustomColumnsFormat + " |CLUSTER_POOL_HOST|NAMESPACE|API_SERVER"
	}
	return nil
}

func (o *Options) validate() error {
	if !helpers.IsOutputFormatSupported(o.OutputFormat) {
		return fmt.Errorf("invalid output format %s", helpers.SupportedOutputFormat)
	}
	return nil
}

func (o *Options) run() (err error) {
	cphs, err := clusterpoolhost.GetClusterPoolHosts()
	if err != nil {
		return err
	}
	helpers.Print(cphs, o.OutputFormat, clusterpoolhost.ConvertClusterPoolHostsForPrint)
	return nil
}
