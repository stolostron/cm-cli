// Copyright Contributors to the Open Cluster Management project

package helpers

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	GvrCC schema.GroupVersionResource = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterclaims"}
	GvrCP schema.GroupVersionResource = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterpools"}
	GvrCD schema.GroupVersionResource = schema.GroupVersionResource{Group: "hive.openshift.io", Version: "v1", Resource: "clusterdeployments"}
)
