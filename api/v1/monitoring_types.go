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

// MonitoringSpec defines the desired state of Monitoring
type MonitoringSpec struct {

	// Important: Run "make" to regenerate code after modifying this file
	// This is where you define what Spec you want for your CR Monitoring
	Source       Source `json:"source"`
	Affix        string `json:"affix,omitempty"`
	DecisionKind string `json:"decisionkind"`
	//MonitoringPolicies MonitoringPolicies `json:"monitoringpolicies"`
	Policy  Policy `json:"policy"`
	Message string `json:"message"`
}

// MonitoringStatus defines the observed state of Monitoring
type MonitoringStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	Affix string `json:"affix,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Monitoring is the Schema for the monitorings API
type Monitoring struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MonitoringSpec   `json:"spec,omitempty"`
	Status MonitoringStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MonitoringList contains a list of Monitoring
type MonitoringList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Monitoring `json:"items"`
}

//-------Additional Structure ---------//

type MonitoringPolicies struct {
	Data          map[string]string `json:"data"`
	TresholdKind  map[string]string `json:"tresholdkind"`
	TresholdValue map[string]string `json:"tresholdvalue"`
	Time          string            `json:"time,omitempty"`
}

type Policy struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Engine      Engine `json:"engine"`
	Rule        Rule   `json:"rule"`
	Kind        string `json:"kind"`
	Data        Data   `json:"data"`
	Input       Input  `json:"input"`
	Result      Result `json:"result"`
}

type Engine struct {
	Kind string `json:"kind"`
	Url  string `json:"url"`
	Api  Api    `json:"api"`
}

type Api struct {
	Policy string `json:"policy"`
	Data   string `json:"data"`
}

type Rule struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

type Data struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

/*type Data struct {
  Name   string `json:"name"`
  Key 	 []string    `json:"key,omitempty"`
  Value  [][]string  `json:"value,omitempty"`

}
*/
type Input struct {
	Schema string `json:"schema"`
	Value  string `json:"value,omitempty"`
}

type Result struct {
	Schema string `json:"schema"`
	Value  string `json:"value,omitempty"`
}

type Source struct {
	Addresse string `json:"addresse"`
	Port     int32  `json:"port"`
	Interval int32  `json:"interval"`
}

func init() {
	SchemeBuilder.Register(&Monitoring{}, &MonitoringList{})
}
