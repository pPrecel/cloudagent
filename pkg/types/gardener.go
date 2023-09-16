package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConditionTrue ConditionStatus = "True"
)

// ShootList is a list of Shoot objects.
type ShootList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list object metadata.
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is the list of Shoots.

	Items []Shoot `json:"items"`
}

// Shoot represents a Shoot cluster created and managed by Gardener.
type Shoot struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Most recently observed status of the Shoot cluster.
	Status ShootStatus `json:"status,omitempty"`
}

// ShootStatus holds the most recently observed status of the Shoot cluster.
type ShootStatus struct {
	// // Conditions represents the latest available observations of a Shoots's current state.
	Conditions []Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`

	// // IsHibernated indicates whether the Shoot is currently hibernated.
	IsHibernated bool `json:"hibernated"`
}

// ConditionStatus is the status of a condition.
type ConditionStatus string

// Condition holds the information about the state of a resource.
type Condition struct {
	// Status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status"`

	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
}
