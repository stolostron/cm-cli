// Copyright Contributors to the Open Cluster Management project
package clusterpool

import (
	"fmt"

	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"

	"github.com/spf13/cobra"
)

const (
	AWS   = "aws"
	AZURE = "azure"
	GCP   = "gcp"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	o.values, err = helpers.ConvertValuesFileToValuesMap(o.valuesPath, "")
	if err != nil {
		return err
	}

	if len(o.values) == 0 {
		return fmt.Errorf("values are missing")
	}

	if len(args) > 0 {
		o.ClusterPool = args[0]
	}

	return nil
}

func (o *Options) validate() (err error) {
	icp, ok := o.values["clusterPool"]
	if !ok || icp == nil {
		return fmt.Errorf("clusterPool is missing")
	}
	cp := icp.(map[string]interface{})
	icloud, ok := cp["cloud"]
	if !ok || icloud == nil {
		return fmt.Errorf("cloud type is missing")
	}
	cloud := icloud.(string)
	if cloud != AWS && cloud != AZURE && cloud != GCP {
		return fmt.Errorf("supported cloud type are (%s, %s, %s) and got %s", AWS, AZURE, GCP, cloud)
	}
	o.cloud = cloud

	_, ocpImageOk := cp["ocpImage"]
	_, imageSetRef := cp["imageSetRef"]
	if ocpImageOk && imageSetRef {
		return fmt.Errorf("ocpImage and imageSetRef are mutually exclusive")
	}

	if o.ClusterPool == "" {
		iname, ok := cp["name"]
		if !ok || iname == nil {
			return fmt.Errorf("clusterPool name is missing")
		}
		o.ClusterPool = iname.(string)
		if len(o.ClusterPool) == 0 {
			return fmt.Errorf("clusterPool.name not specified")
		}
	}

	cp["name"] = o.ClusterPool

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

	err = o.createClusterPool(cphs)

	if len(o.ClusterPoolHost) != 0 {
		if err := cphs.SetActive(currentCph); err != nil {
			return err
		}
	}
	return err
}

func (o *Options) createClusterPool(cphs *clusterpoolhost.ClusterPoolHosts) (err error) {
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
	return clusterpoolhost.CreateClusterPool(o.ClusterPool, o.cloud, o.values, o.CMFlags.DryRun, o.outputFile)
}
