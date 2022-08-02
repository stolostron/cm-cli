// Copyright Contributors to the Open Cluster Management project
package v1alpha1

import (
	"embed"

	"github.com/stolostron/applier/pkg/asset"
)

//go:embed crd
var files embed.FS

func GetScenarioResourcesReader() *asset.ScenarioResourcesReader {
	return asset.NewScenarioResourcesReader(&files)
}
