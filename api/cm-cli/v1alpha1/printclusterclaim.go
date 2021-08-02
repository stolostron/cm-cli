// Copyright Contributors to the Open Cluster Management project

package v1alpha1

import (
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PrintClusterClaimSpec defines the desired state of PrintClusterPool
type PrintClusterClaimSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ClusterPoolHostName string               `json:"clusterPoolHostName"`
	ClusterClaim        *hivev1.ClusterClaim `json:"clusterClaim"`
	Hibernate           string               `json:"hibernate"`
	PowerState          string               `json:"powerState"`
	ID                  string               `json:"id"`
	ErrorMessage        string               `json:"error"`
}

// +kubebuilder:object:root=true

// PrintClusterClaim is the Schema for the authrealms API
type PrintClusterClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PrintClusterClaimSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// PrintClusterClaimList contains a list of PrintClusterPool
type PrintClusterClaimList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// List of PrintClusterPool.
	// +listType=set
	Items []PrintClusterClaim `json:"items"`
}
