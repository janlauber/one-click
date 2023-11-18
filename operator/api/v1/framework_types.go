/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type ImageSpec struct {
	Registry   string `json:"registry"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}

type HorizontalScaleSpec struct {
	MinReplicas                    int32 `json:"minReplicas"`
	MaxReplicas                    int32 `json:"maxReplicas"`
	TargetCPUUtilizationPercentage int32 `json:"targetCPUUtilizationPercentage"`
}

type ResourceRequirements struct {
	Requests ResourceList `json:"requests"`
	Limits   ResourceList `json:"limits"`
}

type ResourceList struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SecretItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type InterfaceSpec struct {
	Name    string        `json:"name"`
	Port    int32         `json:"port"`
	Ingress []IngressSpec `json:"ingress"`
}

type IngressSpec struct {
	Host         string            `json:"host"`
	Path         string            `json:"path"`
	TLS          bool              `json:"tls"`
	IngressClass string            `json:"ingressClass,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty"`
}

// FrameworkSpec defines the desired state of Framework
type FrameworkSpec struct {
	Image              ImageSpec            `json:"image"`
	HorizontalScale    HorizontalScaleSpec  `json:"horizontalScale"`
	Resources          ResourceRequirements `json:"resources"`
	Env                []EnvVar             `json:"env"`
	Secrets            []SecretItem         `json:"secrets"`
	Interfaces         []InterfaceSpec      `json:"interfaces"`
	ServiceAccountName string               `json:"serviceAccountName"`
}

// FrameworkStatus defines the observed state of Framework
type FrameworkStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Framework is the Schema for the frameworks API
type Framework struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FrameworkSpec   `json:"spec,omitempty"`
	Status FrameworkStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FrameworkList contains a list of Framework
type FrameworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Framework `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Framework{}, &FrameworkList{})
}
