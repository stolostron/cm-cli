// Copyright Contributors to the Open Cluster Management project

package helpers

import "k8s.io/apimachinery/pkg/runtime/schema"

const (
	QPS   = 200
	Burst = 200
)

var (
	GvrCC  schema.GroupVersionResource = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterclaims"}
	GvrCP  schema.GroupVersionResource = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterpools"}
	GvrCD  schema.GroupVersionResource = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterdeployments"}
	GvrCIS schema.GroupVersionResource = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterimagesets"}
)
