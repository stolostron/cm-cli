// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"
	"github.com/stolostron/cm-cli/pkg/helpers"
	"k8s.io/kubectl/pkg/cmd/get"
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

	return o.executeCommand(cph)
}

func (o *Options) executeCommand(cph *clusterpoolhost.ClusterPoolHost) (err error) {
	outputFormat := "yaml"
	err = cph.SetClusterClaimContext(o.Cluster, false, o.Timeout, o.CMFlags.DryRun, o.outputFile, &get.PrintFlags{OutputFormat: &outputFormat})
	if err != nil {
		return err
	}
	context := cph.GetClusterContextName(o.Cluster)
	return helpers.ExecuteWithContext(context, os.Args, o.CMFlags.DryRun, o.streams, o.outputFile)
}
