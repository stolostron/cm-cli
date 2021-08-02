// Copyright Contributors to the Open Cluster Management project

package v1alpha1

import (
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PrintClusterPoolSpec defines the desired state of PrintClusterPool
type PrintClusterPoolSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Foo string `json:"foo,omitempty"`
	ClusterPoolHostName string              `json:"clusterPoolHostName"`
	ClusterPool         *hivev1.ClusterPool `json:"clusterPool"`
}

// +kubebuilder:object:root=true

// PrintClusterPool is the Schema for the authrealms API
type PrintClusterPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PrintClusterPoolSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// PrintClusterPoolList contains a list of PrintClusterPool
type PrintClusterPoolList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// List of PrintClusterPool.
	// +listType=set
	Items []PrintClusterPool `json:"items"`
}
