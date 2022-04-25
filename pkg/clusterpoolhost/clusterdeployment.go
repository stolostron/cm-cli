// Copyright Contributors to the Open Cluster Management project
package clusterpoolhost

import (
	"context"
	"fmt"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/stolostron/cm-cli/pkg/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
)

func (cph *ClusterPoolHost) GetClusterDeployment(cc *hivev1.ClusterClaim) (*hivev1.ClusterDeployment, error) {
	clusterPoolRestConfig, err := cph.GetGlobalRestConfig()
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(clusterPoolRestConfig)
	if err != nil {
		return nil, err
	}
	if len(cc.Spec.Namespace) == 0 {
		return nil, fmt.Errorf("something wrong happened, the clusterclaim %s doesn't have a spec.namespace set", cc.Name)
	}
	cdu, err := dynamicClient.Resource(helpers.GvrCD).Namespace(cc.Spec.Namespace).Get(context.TODO(), cc.Spec.Namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cd := &hivev1.ClusterDeployment{}
	if runtime.DefaultUnstructuredConverter.FromUnstructured(cdu.UnstructuredContent(), cd); err != nil {
		return nil, err
	}
	return cd, nil
}
