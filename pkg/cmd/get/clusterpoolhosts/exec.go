// Copyright Contributors to the Open Cluster Management project
package clusterpoolhosts

import (
	printclusterpoolv1alpha1 "github.com/open-cluster-management/cm-cli/api/cm-cli/v1alpha1"
	"github.com/open-cluster-management/cm-cli/pkg/clusterpoolhost"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	if len(*o.PrintFlags.OutputFormat) == 0 {
		o.PrintFlags.OutputFormat = &clusterpoolhost.ClusterPoolHostsColumns
	}
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
	pcphs := clusterpoolhost.ConvertToPrintClusterPoolHostList(cphs)
	pcphs.GetObjectKind().
		SetGroupVersionKind(
			schema.GroupVersionKind{
				Group:   printclusterpoolv1alpha1.GroupName,
				Kind:    "PrintClusterPoolHosts",
				Version: printclusterpoolv1alpha1.GroupVersion.Version})
	helpers.Print(pcphs, o.PrintFlags)
	return nil
}
