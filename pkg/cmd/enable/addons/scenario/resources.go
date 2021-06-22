// Copyright Contributors to the Open Cluster Management project
package scenario

import (
	"embed"

	"open-cluster-management.io/clusteradm/pkg/helpers/asset"
)

//go:embed addons
var files embed.FS

//Note: The other resources are imported in the code by creating a reader on
//  "github.com/open-cluster-management/cm-cli/pkg/cmd/attach/cluster/scenario"
// as we don't want to duplicate yamls

func GetScenarioResourcesReader() *asset.ScenarioResourcesReader {
	return asset.NewScenarioResourcesReader(&files)
}
