package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DNSetSpec struct {
	DNSetBasic     `json:",inline"`
	ScaleStrategy  CloneSetScaleStrategy  `json:"scaleStrategy,omitempty"`
	UpdateStrategy CloneSetUpdateStrategy `json:"updateStrategy,omitempty"`
	CloneSetCommon `json:",inline"`
	// +optional
	Overlay *Overlay `json:"overlay,omitempty"`
}

type DNSetBasic struct {
	PodSet    `json:",inline"`
	LogConfig `json:",inline"`

	// ServiceType is the service type of dn service
	// +optional
	// +kubebuilder:default=ClusterIP
	// +kubebuilder:validation:Enum=ClusterIP;NodePort;LoadBalancer
	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`

	// CacheVolume is the desired local cache volume for DNSet,
	// node storage will be used if not specified
	// +optional
	CacheVolume *Volume `json:"cacheVolume,omitempty"`
}

// TODO: figure out what status should be exposed
type DNSetStatus struct {
	ConditionalStatus `json:",inline"`
}

type DNSetDeps struct {
	LogSetRef `json:",inline"`
}

// +kubebuilder:object:root=true

// A DNSet is a resource that represents a set of MO's DN instances
// +kubebuilder:subresource:status
type DNSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DNSetSpec   `json:"spec,omitempty"`
	Deps   DNSetDeps   `json:"deps,omitempty"`
	Status DNSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DNSetList contains a list of DNSet
type DNSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DNSet{}, &DNSetList{})
}
