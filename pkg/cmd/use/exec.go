// Copyright Contributors to the Open Cluster Management project
package use

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("clustername is missing")
	}
	o.ClusterHostPool = args[0]
	return nil
}

func (o *Options) validate() error {
	return nil
}

func (o *Options) run() (err error) {
	return clusterpoolhost.VerifyContext(o.ClusterHostPool, o.CMFlags.DryRun, o.outputFile)
}
