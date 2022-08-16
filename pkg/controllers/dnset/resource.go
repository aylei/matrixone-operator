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

package dnset

import (
	"bytes"
	"fmt"
	"github.com/matrixorigin/matrixone-operator/api/core/v1alpha1"
	"github.com/matrixorigin/matrixone-operator/pkg/controllers/common"
	"github.com/matrixorigin/matrixone-operator/pkg/utils"
	"github.com/matrixorigin/matrixone-operator/runtime/pkg/util"
	"github.com/openkruise/kruise-api/apps/pub"
	kruise "github.com/openkruise/kruise-api/apps/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func syncReplicas(dn *v1alpha1.DNSet, cs *kruise.CloneSet) {
	cs.Spec.Replicas = &dn.Spec.Replicas
}

func syncPodMeta(dn *v1alpha1.DNSet, cs *kruise.CloneSet) {
	dn.Spec.Overlay.OverlayPodMeta(&cs.Spec.Template.ObjectMeta)
}

func syncPodSpec(dn *v1alpha1.DNSet, cs *kruise.CloneSet) {
	main := corev1.Container{
		Name:      v1alpha1.ContainerMain,
		Image:     dn.Spec.Image,
		Resources: dn.Spec.Resources,
		Command: []string{
			"/bin/sh", fmt.Sprintf("%s/%s", common.ConfigPath, Entrypoint),
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: common.DataVolume, ReadOnly: true, MountPath: common.DataPath},
			{Name: common.ConfigVolume, ReadOnly: true, MountPath: common.ConfigPath},
		},
		Env: []corev1.EnvVar{
			util.FieldRefEnv(common.PodNameEnvKey, "metadata.name"),
			util.FieldRefEnv(common.NamespaceEnvKey, "metadata.namespace"),
			util.FieldRefEnv(common.PodIPEnvKey, "status.podIP"),
			{Name: common.HeadlessSvcEnvKey, Value: utils.GetHeadlessSvcName(dn)},
		},
	}
	dn.Spec.Overlay.OverlayMainContainer(&main)
	podSpec := corev1.PodSpec{
		Containers: []corev1.Container{main},
		ReadinessGates: []corev1.PodReadinessGate{{
			ConditionType: pub.InPlaceUpdateReady,
		}},
	}
	common.SyncTopology(dn.Spec.TopologyEvenSpread, &podSpec)

	dn.Spec.Overlay.OverlayPodSpec(&podSpec)
	cs.Spec.Template.Spec = podSpec
}

// buildDNSetConfigMap return dn set configmap
func buildDNSetConfigMap(dn *v1alpha1.DNSet) (*corev1.ConfigMap, error) {

	buff := new(bytes.Buffer)
	//err := startScriptTpl.Execute(buff, &model{
	//	ConfigFilePath: fmt.Sprintf("%s/%s", configPath, ConfigFile),
	//})
	//if err != nil {
	//	panic(err)
	//}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: utils.GetNamespace(dn),
			Name:      utils.GetConfigName(dn),
			Labels:    common.SubResourceLabels(dn),
		},
		Data: map[string]string{
			//ConfigFile: s,
			Entrypoint: buff.String(),
		},
	}, nil
}

func buildHeadlessSvc(dn *v1alpha1.DNSet) *corev1.Service {
	ports := getDNServicePort()
	return common.GetHeadlessService(dn, ports)
}

func buildDiscoverySvc(dn *v1alpha1.DNSet) *corev1.Service {
	ports := getDNServicePort()
	return common.GetDiscoveryService(dn, ports, dn.Spec.ServiceType)
}

func buildDNSet(dn *v1alpha1.DNSet) *kruise.CloneSet {
	return common.GetCloneSet(dn)
}

func syncPersistentVolumeClaim(dn *v1alpha1.DNSet, cloneSet *kruise.CloneSet) {
	if dn.Spec.CacheVolume != nil {
		dataPVC := common.GetPersistentVolumeClaim(dn.Spec.CacheVolume.Size, dn.Spec.CacheVolume.StorageClassName)
		tpls := []corev1.PersistentVolumeClaim{dataPVC}
		dn.Spec.Overlay.AppendVolumeClaims(&tpls)
		cloneSet.Spec.VolumeClaimTemplates = tpls
	}
}
