// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"github.com/open-cluster-management/cm-cli/pkg/genericclioptions"
	"github.com/open-cluster-management/cm-cli/pkg/helpers"
)

func IsSupported(cmFlags *genericclioptions.CMFlags) (isSupported bool, err error) {
	var serverNamespace string
	switch {
	case len(cmFlags.ServerNamespace) != 0:
		serverNamespace = cmFlags.ServerNamespace
	default:
		cph, err := GetCurrentClusterPoolHost()
		if err != nil {
			return false, err
		}
		serverNamespace = cph.ServerNamespace
	}
	cmFlags.ServerNamespace = serverNamespace
	return helpers.IsSupported(cmFlags)
}
