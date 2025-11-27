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

type Mode string

const (
	ONE_SHOT    Mode = "OneShot"
	ONCE_BY_TAG Mode = "OnceByTag"
	RECURRENT   Mode = "Recurrent"
)

// ImageSpec defines the desired state of Image
type ImageSpec struct {
	AllowCandidateRelease bool `json:"allowCandidateRelease,omitempty"`
	// +kubebuilder:validation:Required
	Destination ImageEndpoint `json:"destination,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="5m"
	Frequency string `json:"frequency,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=OneShot;OnceByTag;Recurrent
	// Select mode you want to apply to the job
	Mode Mode `json:"mode,omitempty"`
	// +kubebuilder:validation:Required
	Source ImageEndpoint `json:"source,omitempty"`
	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// +kubebuilder:validation:Optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
}

type History struct {
	PerformedAt metav1.Time `json:"performedAt,omitempty"`
}

type ImageStatus struct {
	History []History `json:"history,omitempty"`
	// +kubebuilder:default:=COMPLETED
	Phase            string   `json:"phase,omitempty"`
	TagAlreadySynced []string `json:"tagAlreadySynced,omitempty"`
	// +kubebuilder:default:=0
	LastGenerationSeen int64 `json:"lastGenerationSeen,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced

// Image is the Schema for the images API
type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSpec   `json:"spec,omitempty"`
	Status ImageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ImageList contains a list of Image
type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Image{}, &ImageList{})
}
