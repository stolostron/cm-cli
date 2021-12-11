// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
)

func IsSupported(cmFlags *genericclioptions.CMFlags) (isSupported bool, err error) {
	var productNamespace string
	switch {
	case len(cmFlags.ProductNamespace) != 0:
		productNamespace = cmFlags.ProductNamespace
	default:
		cph, err := GetCurrentClusterPoolHost()
		if err != nil {
			return false, err
		}
		productNamespace = cph.ProductNamespace
	}
	cmFlags.ProductNamespace = productNamespace
	return helpers.IsRHACM(cmFlags) || helpers.IsMCE(cmFlags), err
}
