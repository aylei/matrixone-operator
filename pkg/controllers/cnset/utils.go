// Copyright 2022 Matrix Origin
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

package cnset

import (
	"github.com/matrixorigin/matrixone-operator/api/core/v1alpha1"
	kruise "github.com/openkruise/kruise-api/apps/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func getCNServicePort(cn *v1alpha1.CNSet) []corev1.ServicePort {
	return []corev1.ServicePort{
		{
			Name: "service",
			Port: servicePort,
		},
	}
}

func getScaleStrategyConfig(cn *v1alpha1.CNSet) kruise.CloneSetScaleStrategy {
	return kruise.CloneSetScaleStrategy{
		PodsToDelete:   cn.Spec.ScaleStrategy.PodsToDelete,
		MaxUnavailable: cn.Spec.ScaleStrategy.MaxUnavailable,
	}
}

func getUpdateStrategyConfig(cn *v1alpha1.CNSet) kruise.CloneSetUpdateStrategy {
	return kruise.CloneSetUpdateStrategy{
		Type:                  cn.Spec.UpdateStrategy.Type,
		Partition:             cn.Spec.UpdateStrategy.Partition,
		MaxUnavailable:        cn.Spec.UpdateStrategy.MaxUnavailable,
		MaxSurge:              cn.Spec.UpdateStrategy.MaxSurge,
		Paused:                cn.Spec.UpdateStrategy.Paused,
		PriorityStrategy:      cn.Spec.UpdateStrategy.PriorityStrategy,
		ScatterStrategy:       cn.Spec.UpdateStrategy.ScatterStrategy,
		InPlaceUpdateStrategy: cn.Spec.UpdateStrategy.InPlaceUpdateStrategy,
	}
}
