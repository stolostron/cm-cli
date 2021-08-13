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

	cph, err := cphs.GetClusterPoolHostOrCurrent(o.ClusterPoolHost)
	if err != nil {
		return err
	}

	return cph.CreateClusterPool(o.ClusterPool, o.cloud, o.values, o.CMFlags.DryRun, o.outputFile)
}
