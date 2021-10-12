// Copyright Contributors to the Open Cluster Management project

package v1alpha1

import (
	policyv1 "github.com/open-cluster-management/governance-policy-propagator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PrintPoliciesSpec defines the desired state of PrintPolicies
type PrintPoliciesSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS
	// Important: Run "make" to regenerate code after modifying this file

	Policy policyv1.Policy `json:"policy"`
	Age    string          `json:"age"`
}

// PrintPolicies is the Schema for the authrealms API
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=printpolicies
// +kubebuilder:printcolumn:name="Policy Name",type="string",JSONPath=".metadata.name"
// +kubebuilder:printcolumn:name="Namespace",type="string",JSONPath=".metadata.namespace"
// +kubebuilder:printcolumn:name="Compliance State",type="string",JSONPath=".spec.policy.status.compliant"
// +kubebuilder:printcolumn:name="Remediation Action",type="string",JSONPath=".spec.policy.spec.remediationAction"
// +kubebuilder:printcolumn:name="Disabled",type="bool",JSONPath=".spec.policy.spec.disabled"
// +kubebuilder:printcolumn:name="Standards",type="string",JSONPath=".spec.policy.metadata.annotations.policy\\.open-cluster-management\\.io/standards",priority=1
// +kubebuilder:printcolumn:name="Categories",type="string",JSONPath=".spec.policy.metadata.annotations.policy\\.open-cluster-management\\.io/categories",priority=1
// +kubebuilder:printcolumn:name="Controls",type="string",JSONPath=".spec.policy.metadata.annotations.policy\\.open-cluster-management\\.io/controls",priority=1
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".spec.age"

type PrintPolicies struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PrintPoliciesSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// PrintPoliciesList contains a list of PrintPolicies
type PrintPoliciesList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// List of PrintPolicies.
	// +listType=set
	Items []PrintPolicies `json:"items"`
}
