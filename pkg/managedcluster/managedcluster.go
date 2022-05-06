// Copyright Contributors to the Open Cluster Management project
package managedcluster

import (
	"fmt"
	"strconv"

	"github.com/stolostron/cm-cli/pkg/helpers"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
)

const (
	ConsoleURLClusterClaim    string = "consoleurl.cluster.open-cluster-management.io"
	HostedClusterClusterClaim string = "hostedcluster.hypershift.openshift.io"

	HostedType       string = "hosted"
	ClusterClaimType string = "clusterclaim"
)

func GetConsoleURL(mc *clusterv1.ManagedCluster) (string, error) {
	for _, c := range mc.Status.ClusterClaims {
		if c.Name == ConsoleURLClusterClaim {
			return c.Value, nil
		}
	}
	return "", fmt.Errorf("console url not found for cluster %s", mc.Name)

}

func IsHosted(mc *clusterv1.ManagedCluster) bool {
	for _, c := range mc.Status.ClusterClaims {
		if c.Name == HostedClusterClusterClaim {
			b, err := strconv.ParseBool(c.Value)
			if err != nil {
				return false
			}
			return b
		}
	}
	return false
}

func GetClusterType(mc *clusterv1.ManagedCluster) string {
	if IsHosted(mc) {
		return HostedType
	}
	return ClusterClaimType
}

func OpenManagedCluster(mc *clusterv1.ManagedCluster) error {

	consoleURL, err := GetConsoleURL(mc)
	if err != nil {
		return err
	}

	return helpers.Openbrowser(consoleURL)
}
