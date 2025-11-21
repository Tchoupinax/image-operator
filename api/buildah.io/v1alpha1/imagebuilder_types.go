/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ImageEndpoint struct {
	// +kubebuilder:validation:Required
	ImageName string `json:"name,omitempty"`
	// +kubebuilder:validation:Required
	ImageVersion string `json:"version,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	// With EKS you might want to use IRSA (Iam Roles for Service Accounts)
	// In this case, additionnal operation have to be done in the job.
	// This option cares of it
	UseAwsIRSA bool `json:"useAwsIRSA,omitempty"`
}

type Resource struct {
	// +kubebuilder:validation:Optional
	Memory string `json:"memory,omitempty"`
	// +kubebuilder:validation:Optional
	Cpu string `json:"cpu,omitempty"`
}

type ResourcesQuota struct {
	// +kubebuilder:validation:Optional
	Limits Resource `json:"limits,omitempty"`
	// +kubebuilder:validation:Optional
	Requests Resource `json:"requests,omitempty"`
}

type Architecture string

const (
	ARM64 Architecture = "arm64"
	AMD64 Architecture = "amd64"
	Both  Architecture = "both"
)

// ImageBuilderSpec defines the desired state of ImageBuilder
type ImageBuilderSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=arm64;amd64;both
	// Select mode you want to apply to the job
	Architecture Architecture `json:"architecture,omitempty"`
	// +kubebuilder:validation:Required
	Image ImageEndpoint `json:"image,omitempty"`
	// +kubebuilder:validation:Required
	Source string `json:"source,omitempty"`
	// +kubebuilder:validation:Optional
	Resources ResourcesQuota `json:"resources,omitempty"`
	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// +kubebuilder:validation:Optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
}

// ImageBuilderStatus defines the observed state of ImageBuilder
type ImageBuilderStatus struct {
	RanJobs []string `json:"ranJobs,omitempty"`
	Pushed  bool     `json:"pushed,omitempty"`
	// +kubebuilder:default:=0
	LastGenerationSeen int64 `json:"lastGenerationSeen,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ImageBuilder is the Schema for the imagebuilders API
type ImageBuilder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageBuilderSpec   `json:"spec,omitempty"`
	Status ImageBuilderStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ImageBuilderList contains a list of ImageBuilder
type ImageBuilderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImageBuilder `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ImageBuilder{}, &ImageBuilderList{})
}
