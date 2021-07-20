// Copyright Contributors to the Open Cluster Management project
package clusterpoolhosts

import (
	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
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
	if o.raw {
		return cphs.RawPrint()
	}
	cphs.Print()
	return nil
}
