// Copyright Contributors to the Open Cluster Management project
package clusterclaim

import (
	"fmt"
	"os"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/spf13/cobra"
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
	cphs, err := clusterpoolhost.GetClusterPoolHosts()
	if err != nil {
		return err
	}

	currentCph, err := cphs.GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}

	if len(o.ClusterPoolHost) != 0 {
		cph, err := cphs.GetClusterPoolHost(o.ClusterPoolHost)
		if err != nil {
			return err
		}

		err = cphs.SetActive(cph)
		if err != nil {
			return err
		}
	}
	err = o.executeCommand(cphs)

	if len(o.ClusterPoolHost) != 0 {
		if err := cphs.SetActive(currentCph); err != nil {
			return err
		}
	}
	return err
}

func (o *Options) executeCommand(cphs *clusterpoolhost.ClusterPoolHosts) (err error) {
	if len(o.ClusterPoolHost) != 0 {
		cph, err := cphs.GetClusterPoolHost(o.ClusterPoolHost)
		if err != nil {
			return err
		}

		err = cphs.SetActive(cph)
		if err != nil {
			return err
		}
	}
	cph, err := cphs.GetCurrentClusterPoolHost()
	if err != nil {
		return err
	}
	err = cph.SetClusterClaimContext(o.Cluster, false, o.Timeout, o.CMFlags.DryRun, o.outputFile)
	if err != nil {
		return err
	}
	context := cph.GetClusterContextName(o.Cluster)
	return helpers.ExecuteWithContext(context, os.Args, o.CMFlags.DryRun, o.streams, o.outputFile)
}
