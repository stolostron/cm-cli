// Copyright Contributors to the Open Cluster Management project

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PrintClusterPoolHostSpec defines the desired state of PrintClusterPool
type PrintClusterPoolHostSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Name of the cluster pool
	Name string `json:"name"`
	// true if this cluster pool is the Active one
	Active bool `json:"active"`
	// The API address of the cluster where your `ClusterPools` are defined. Also referred to as the "ClusterPool host"
	APIServer string `json:"apiServer"`
	// The URL of the OpenShift console for the ClusterPool host
	Console string `json:"console"`
	// Namespace where `ClusterPools` are defined
	Namespace string `json:"namespace"`
	// Name of a `Group` (`user.openshift.io/v1`) that should be added to each `ClusterClaim` for team access
	Group string `json:"group"`
}

// PrintClusterPoolHost is the Schema for the authrealms API
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=printclusterpoolhosts
// +kubebuilder:printcolumn:name="Name",type="string",JSONPath=".metadata.name"
// +kubebuilder:printcolumn:name="Active",type="string",JSONPath=".spec.active"
// +kubebuilder:printcolumn:name="Namespace",type="string",JSONPath=".spec.Namespace"
// +kubebuilder:printcolumn:name="Api_server",type="string",JSONPath=".spec.APIServer"
// +kubebuilder:printcolumn:name="Console",type="string",JSONPath=".spec.Console"

type PrintClusterPoolHost struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PrintClusterPoolHostSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// PrintClusterPoolHostList contains a list of PrintClusterPool
type PrintClusterPoolHostList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// List of PrintClusterPool.
	// +listType=set
	Items []PrintClusterPoolHost `json:"items"`
}
