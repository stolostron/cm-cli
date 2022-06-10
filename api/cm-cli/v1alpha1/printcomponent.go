// Copyright Contributors to the Open Cluster Management project

package v1alpha1

import (
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PrintComponentSpec defines the desired state of PrintClusterPool
type PrintComponentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Enabled      bool                 `json:"enabled"`
	ClusterClaim *hivev1.ClusterClaim `json:"clusterClaim"`
}

// PrintComponents the Schema for the authrealms API
//	ComponentsColumns            string = "custom-columns=CLUSTER_POOL_HOST:.spec.clusterPoolHostName,CLUSTER_CLAIM:.spec.clusterClaim.Name,POWER_STATE:.spec.powerState,HIBERNATE:.spec.hibernate,ID:.spec.ID,ERROR:.spec.errorMessage"
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=printcomponents
// +kubebuilder:printcolumn:name="NAME",type="string",JSONPath=".metadata.name"
// +kubebuilder:printcolumn:name="ENABLED",type="string",JSONPath=".spec.enabled"

type PrintComponent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PrintComponentSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// PrintComponentList contains a list of PrintClusterPool
type PrintComponentList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// List of PrintClusterPool.
	// +listType=set
	Items []PrintComponent `json:"items"`
}
