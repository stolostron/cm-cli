// Copyright Contributors to the Open Cluster Management project

package version

import (
	_ "embed"
)

//go:embed VERSION.txt
var version []byte

func GetVersion() string {
	return string(version)
}
