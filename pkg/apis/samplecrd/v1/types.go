package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SampleCRD is a specification for a SampleCRD resource
type SampleCRD struct {
	metav1.TypeMeta     `json:",inline"`
	metav1.ObjectMeta   `json:"metadata,omitempty"`

	Spec SampleSpec     `json:"spec"`
	Status SampleStatus `json:"status"`
}

type SampleSpec struct {
	DeploymentName string `json:"deploymentName"`
	Replicas *int32       `json:"replicas"`
	Pods []v1.Pod         `json:"pods"`
}

type SampleStatus struct {
	CurrentReplicas int32 `json:"currentReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SampleCRDList is a list of SampleCRD resources
type SampleCRDList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []SampleCRD `json:"items"`
}