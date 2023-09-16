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

type Svc struct {
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	Port       int32  `json:"port,omitempty"`
	TargetPort int    `json:"targetPort,omitempty"`
	NodePort   int32  `json:"nodePort,omitempty"`
}
type Module struct {
	Name  string            `json:"name,omitempty"`
	Image string            `json:"image,omitempty"`
	Port  int32             `json:"port,omitempty"`
	Env   map[string]string `json:"env,omitempty"`
	Svc   Svc               `json:"svc,omitempty"`
}

type Volume struct {
	Capacity string `json:"capacity,omitempty"`
	Path     string `json:"path,omitempty"`
	Storage  string `json:"storage,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CrudSpec defines the desired state of Crud
type CrudSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	App    Module `json:"app,omitempty"`
	Db     Module `json:"db,omitempty"`
	Volume Volume `json:"volume,omitempty"`
}

// CrudStatus defines the observed state of Crud
type CrudStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Crud is the Schema for the cruds API
type Crud struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CrudSpec   `json:"spec,omitempty"`
	Status CrudStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CrudList contains a list of Crud
type CrudList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Crud `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Crud{}, &CrudList{})
}
