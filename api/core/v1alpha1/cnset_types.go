package v1alpha1

import (
	recon "github.com/matrixorigin/matrixone-operator/runtime/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type CNSetSpec struct {
	CNSetBasic `json:",inline"`

	// +optional
	Overlay *Overlay `json:"overlay,omitempty"`
}

type CNSetBasic struct {
	PodSet `json:",inline"`

	// CloneSetCommon kruise clone set common config for cnset
	// +optional
	CloneSetCommon `json:",inline"`

	// ScaleStrategy clone set scalestrategy config for cn service
	// +optional
	ScaleStrategy CloneSetScaleStrategy `json:"scaleStrategy,omitempty"`

	// Updatestrategy clone set updatestrategy config for cn service
	// +optional
	UpdateStrategy CloneSetUpdateStrategy `json:"updateStrategy,omitempty"`

	InitialConfig CNInitialConfig `json:"initialConfig,omitempty"`

	// ConfigMap is reference to a key in a config map
	// +optional
	ConfigMap *corev1.ConfigMapKeySelector `json:"configmap,omitempty"`

	// ServiceType is the service type of cn service
	// +optional
	// +kubebuilder:default=ClusterIP
	// +kubebuilder:validation:Enum=ClusterIP;NodePort;LoadBalancer
	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`

	// CacheVolume is the desired local cache volume for CNSet,
	// node storage will be used if not specified
	// +optional
	CacheVolume *Volume `json:"cacheVolume,omitempty"`
}

type CNInitialConfig struct {
	// TPAvailable maximum ap set is available
	// +optional
	// kubebuilder:default=0
	APAvailable *intstr.IntOrString `json:"apAvailable,omitempty"`
}

// TODO: figure out what status should be exposed
type CNSetStatus struct {
	ConditionalStatus `json:",inline"`
}

type CNSetDeps struct {
	LogSetRef `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// A CNSet is a resource that represents a set of MO's CN instances
// +kubebuilder:subresource:status
type CNSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CNSetSpec   `json:"spec,omitempty"`
	Deps   CNSetDeps   `json:"deps,omitempty"`
	Status CNSetStatus `json:"status,omitempty"`
}

func (d *CNSet) GetDependencies() []recon.Dependency {
	var deps []recon.Dependency
	if d.Deps.LogSet != nil {
		deps = append(deps, &recon.ObjectDependency[*LogSet]{
			ObjectRef: d.Deps.LogSet,
			ReadyFunc: func(l *LogSet) bool {
				return l.Status.Ready()
			},
		})
	}
	return deps
}

//+kubebuilder:object:root=true

// CNSetList contains a list of CNSet
type CNSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CNSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CNSet{}, &CNSetList{})
}
