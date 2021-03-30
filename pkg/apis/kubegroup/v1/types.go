package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 添加`+genclient:noStatus`标记可以不添加自定义资源的Status成员.
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodCluster describes a PodCluster resource
type PodCluster struct {
	// TypeMeta为各资源通用元信息, 包括kind和apiVersion.
	metav1.TypeMeta `json:",inline"`
	// ObjectMeta为特定类型的元信息, 包括name, namespace, selfLink, labels等.
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// spec字段
	Spec PodClusterSpec `json:"spec"`
	// status字段
	Status PodClusterStatus `json:"status"`
}

// PodClusterSpec is the spec for a MyResource resource
type PodClusterSpec struct {
	PodReplicas int32 `json:"PodReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodClusterList is a list of PodCluster resources
type PodClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []PodCluster `json:"items"`
}

// PodClusterStatus is the status for a PodStatus resource
type PodClusterStatus struct {
	PodReplicas int32 `json:"PodReplicas"`
}
