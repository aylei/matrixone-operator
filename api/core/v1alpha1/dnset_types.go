package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DNSetSpec struct {
	PodSet `json:",inline"`

	// Volume is the local persistent volume for each dn service isntance
	// +required
	Volume Volume `json:"volume"`

	// ShardStorage is an external shard storage shared by all dn service instances
	// +required
	SharedStorage SharedStorageProvider `json:"sharedStorage"`

	// CacheVolume is the desired local cache volume for DNSet,
	// node storage will be used if not specified
	// +optional
	CacheVolume *Volume `json:"cacheVolume,omitempty"`
}

// TODO: figure out what status should be exposed
type DNSetStatus struct {
	ConditionalStatus `json:",inline"`

	AvailableStores []DNStore `json:"availableStores,omitempty"`
	FailedStores    []DNStore `json:"failedStores,omitempty"`
}

type DNStore struct {
	PodName           string `json:"podName,omitempty"`
	Phase             string `json:"phase,omitempty"`
	LastTrasitionTime string `json:"lastTransitionTime,omitempty"`
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
