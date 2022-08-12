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

package matrixone

import (
	"github.com/matrixorigin/matrixone-operator/api/core/v1alpha1"
	recon "github.com/matrixorigin/matrixone-operator/runtime/pkg/reconciler"
	"k8s.io/klog/v2"
)

func syncDNBasic(mo *v1alpha1.MatrixOneCluster, dn *v1alpha1.DNSet) error {

	klog.V(recon.Debug).Info("sync DN basic...")
	dn.Name = mo.Name + "-dn"
	dn.Namespace = mo.Namespace
	dn.Spec.DNSetBasic = mo.Spec.DN
	dn.Kind = "DNSet"

	return nil
}

func syncCNBasic(mo *v1alpha1.MatrixOneCluster, cn *v1alpha1.CNSet) error {

	klog.V(recon.Debug).Info("sync CN basic...")
	cn.Name = mo.Name + "-cn"
	cn.Namespace = mo.Namespace
	cn.Spec.CNSetBasic = mo.Spec.CN
	cn.Kind = "CNSet"

	return nil
}

func syncLogServiceBasic(mo *v1alpha1.MatrixOneCluster, logService *v1alpha1.LogSet) error {

	klog.V(recon.Debug).Info("sync logService...")

	logService.Name = mo.Name + "-log"
	logService.Namespace = mo.Namespace
	logService.Kind = "LogSet"

	return nil
}
