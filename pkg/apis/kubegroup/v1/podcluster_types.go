


package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodCluster
// +k8s:openapi-gen=true
// +resource:path=podclusters,strategy=PodClusterStrategy
type PodCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodClusterSpec   `json:"spec,omitempty"`
	Status PodClusterStatus `json:"status,omitempty"`
}

// PodClusterSpec defines the desired state of PodCluster
type PodClusterSpec struct {
	PodReplicas int32 `json:"podReplicas"`
}

// PodClusterStatus defines the observed state of PodCluster
type PodClusterStatus struct {
}
