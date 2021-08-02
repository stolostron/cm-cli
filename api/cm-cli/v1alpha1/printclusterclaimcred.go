// Copyright Contributors to the Open Cluster Management project

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PrintClusterClaimCredentialSpec defines the desired state of PrintClusterPool
type PrintClusterClaimCredentialSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	User       string `json:"user"`
	Password   string `json:"pasword"`
	Basedomain string `json:"baseDomain"`
	ApiUrl     string `json:"apiServer"`
	ConsoleUrl string `json:"console"`
}

// +kubebuilder:object:root=true

// PrintClusterClaimCredential is the Schema for the authrealms API
type PrintClusterClaimCredential struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PrintClusterClaimCredentialSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// PrintClusterClaimCredentialList contains a list of PrintClusterPool
type PrintClusterClaimCredentialList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// List of PrintClusterPool.
	// +listType=set
	Items []PrintClusterClaimCredential `json:"items"`
}
