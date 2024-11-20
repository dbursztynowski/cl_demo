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

// ExecutionSpec defines the desired state of Execution
type ExecutionSpec struct {
	// Important: Run "make" to regenerate code after modifying this file
	// This is where you define what Spec you want for your CR Decision
	Affix           string `json:"affix,omitempty"`
	Time            string `json:"time,omitempty"`
	Message         string `json:"message"`
	ExecutionTypeId int32  `json:"executiontypeid"`
	Config          string `json:"config"`
}

// ExecutionStatus defines the observed state of Execution
type ExecutionStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
	Affix string `json:"affix,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Execution is the Schema for the executions API
type Execution struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExecutionSpec   `json:"spec,omitempty"`
	Status ExecutionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ExecutionList contains a list of Execution
type ExecutionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Execution `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Execution{}, &ExecutionList{})
}
