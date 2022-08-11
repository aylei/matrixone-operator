// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MatrixOneClusterSpec defines the desired state of MatrixOneCluster
// Note that MatrixOneCluster does not support specify overlay for underlying sets directly due to the size limitation
// of kubernetes apiserver
type MatrixOneClusterSpec struct {
	//  CN is the default CN pod set of this Cluster
	// +required
	CN CNSetBasic `json:"cn"`

	// DN is the default DN pod set of this Cluster
	DN DNSetBasic `json:"dn"`

	// LogService is the default LogService pod set of this cluster
	LogService LogSetBasic `json:"logService"`

	// Version is the version of the cluster, which translated
	// to the docker image tag used for each component.
	// default to the recommended version of the operator
	// +optional
	Version *string `json:"version"`

	// WebUIEnabled indicates whether deploy the MO web-ui,
	// default to true.
	// +optional
	WebUIEnabled *bool `json:"webUIEnabled,omitempty"`

	// ImageRepository allows user to override the default image
	// repository in order to use a docker registry proxy or private
	// registry.
	// +optional
	ImageRepository *string `json:"imageRepository,omitempty"`
}

// MatrixOneClusterStatus defines the observed state of MatrixOneCluster
type MatrixOneClusterStatus struct {
	ConditionalStatus `json:",inline"`
	// CN is the CN set status
	CN *CNSetStatus `json:"cn,omitempty"`
	// DN is the DN set status
	DN *DNSetStatus `json:"dn,omitempty"`
	// LogService is the LogService status
	LogService *LogSetStatus `json:"logService,omitempty"`
}

// +kubebuilder:object:root=true

// A MatrixOneCluster is a resource that represents a MatrixOne Cluster
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=mo
type MatrixOneCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MatrixOneClusterSpec   `json:"spec,omitempty"`
	Status MatrixOneClusterStatus `json:"status,omitempty"`
}

// MatrixOneGlobalSpec Global config for MatrixOne, like log etc.
//type MatrixOneGlobalSpec struct {
//	LogConfig `json:",inline"`
//	Image     string `json:"image,omitempty"`
//}

//+kubebuilder:object:root=true

// MatrixOneClusterList contains a list of MatrixOneCluster
type MatrixOneClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MatrixOneCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MatrixOneCluster{}, &MatrixOneClusterList{})
}
