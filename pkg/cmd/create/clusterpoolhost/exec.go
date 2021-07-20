// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"fmt"
	"net/url"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"

	"github.com/spf13/cobra"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("clustername is missing")
	}
	o.ClusterPoolHost.Name = args[0]
	return nil
}

func (o *Options) validate() error {
	if len(o.ClusterPoolHost.APIServer) == 0 {
		return fmt.Errorf("api-server is missing")
	}
	_, err := url.Parse(o.ClusterPoolHost.APIServer)
	if err != nil {
		return err
	}
	if len(o.ClusterPoolHost.Console) == 0 {
		return fmt.Errorf("console is missing")
	}
	_, err = url.Parse(o.ClusterPoolHost.Console)
	if err != nil {
		return err
	}
	if len(o.ClusterPoolHost.Group) == 0 {
		return fmt.Errorf("group is missing")
	}
	if len(o.ClusterPoolHost.Namespace) == 0 {
		return fmt.Errorf("namespace is missing")
	}
	_, err = clusterpoolhost.GetClusterPoolHost(o.ClusterPoolHost.Name)
	if err == nil && !o.force {
		return fmt.Errorf("clusterpoolhost already exists, use --force to overwrite")
	}

	return nil
}

func (o *Options) run() (err error) {
	return o.initclusterpoolhost()
}

func (o *Options) initclusterpoolhost() error {
	cph := &clusterpoolhost.ClusterPoolHost{
		Name:      o.ClusterPoolHost.Name,
		APIServer: o.ClusterPoolHost.APIServer,
		Console:   o.ClusterPoolHost.Console,
		Group:     o.ClusterPoolHost.Group,
		Namespace: o.ClusterPoolHost.Namespace,
	}
	err := cph.VerifyClusterPoolContext(o.CMFlags.DryRun, o.outputFile)
	if err != nil {
		return err
	}
	cphs, err := clusterpoolhost.GetClusterPoolHosts()
	if err != nil {
		return err
	}

	cph, err = cphs.GetClusterPoolHost(o.ClusterPoolHost.Name)
	if err != nil {
		return err
	}

	return cphs.SetActive(cph)
}
