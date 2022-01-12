// Copyright Contributors to the Open Cluster Management project
package clusterpoolhosts

import (
	"github.com/spf13/cobra"
	printclusterpoolv1alpha1 "github.com/stolostron/cm-cli/api/cm-cli/v1alpha1"
	"github.com/stolostron/cm-cli/pkg/clusterpoolhost"
	"github.com/stolostron/cm-cli/pkg/helpers"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	pcphs := clusterpoolhost.ConvertToPrintClusterPoolHostList(cphs)
	pcphs.GetObjectKind().
		SetGroupVersionKind(
			schema.GroupVersionKind{
				Group:   printclusterpoolv1alpha1.GroupName,
				Kind:    "PrintClusterPoolHost",
				Version: printclusterpoolv1alpha1.GroupVersion.Version})
	// sort.Sort(pcphs.Items)
	helpers.Print(pcphs, o.GetOptions.PrintFlags)
	return nil
}
