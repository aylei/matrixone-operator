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

package common

import (
	"github.com/matrixorigin/matrixone-operator/api/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	SecretKeyAWSS3AccessKeyID     string = "AccessKeyID"
	SecretKeyAWSS3SecretAccessKey string = "SecretAccessKey"
	SecretKeyAWSS3Region          string = "Region"
)

type MatrixOneObject struct {
	scheme *runtime.Scheme
	mo     *v1alpha1.MatrixOneCluster
	cn     *v1alpha1.CNSet
	dn     *v1alpha1.DNSet
}

func (m *MatrixOneObject) setEnvsForS3(secretRef string) []corev1.EnvVar {
	secret := corev1.LocalObjectReference{
		Name: secretRef,
	}
	return []corev1.EnvVar{
		{
			Name: "AWS_REGION",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: secret,
				},
			},
		},
		{
			Name: "AWS_ACCESS_KEY_ID",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: secret,
					Key:                  SecretKeyAWSS3AccessKeyID,
				},
			},
		},
		{
			Name: "AWS_SECRET_ACCESS_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: secret,
					Key:                  SecretKeyAWSS3SecretAccessKey,
				},
			},
		},
	}
}
