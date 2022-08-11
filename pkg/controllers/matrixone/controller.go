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
	"github.com/matrixorigin/matrixone-operator/pkg/controllers/cnset"
	"github.com/matrixorigin/matrixone-operator/pkg/controllers/dnset"
	"github.com/matrixorigin/matrixone-operator/pkg/controllers/logset"
	recon "github.com/matrixorigin/matrixone-operator/runtime/pkg/reconciler"
	"github.com/matrixorigin/matrixone-operator/runtime/pkg/util"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MatrixOneActor struct {
	DNSet  dnset.DNSetActor
	LogSet logset.LogSetActor
	CNSet  cnset.CNSetActor
}

var _ recon.Actor[*v1alpha1.MatrixOneCluster] = &MatrixOneActor{}

func (m *MatrixOneActor) Finalize(ctx *recon.Context[*v1alpha1.MatrixOneCluster]) (bool, error) {
	mo := ctx.Obj

	ctx.Log.Info("finalzie matrixone")

	objs := []client.Object{&v1alpha1.LogSet{}, &v1alpha1.DNSet{}, &v1alpha1.CNSet{}}
	for _, obj := range objs {
		obj.SetNamespace(mo.Namespace)
		if err := util.Ignore(apierrors.IsNotFound, ctx.Delete(obj)); err != nil {
			return false, err
		}
	}
	for _, obj := range objs {
		exist, err := ctx.Exist(client.ObjectKeyFromObject(obj), obj)
		if err != nil {
			return false, err
		}
		if exist {
			return false, nil
		}
	}
	return true, nil

}

func (m *MatrixOneActor) Observe(
	ctx *recon.Context[*v1alpha1.MatrixOneCluster]) (recon.Action[*v1alpha1.MatrixOneCluster], error) {
	mo := ctx.Obj

	logSet := &v1alpha1.LogSet{}
	err, foundLogSet := util.IsFound(ctx.Get(client.ObjectKey{Namespace: mo.Namespace, Name: mo.Name}, logSet))
	if err != nil {
		return nil, errors.Wrap(err, "get log service")
	}

	dnSet := &v1alpha1.DNSet{}
	err, foundDNSet := util.IsFound(ctx.Get(client.ObjectKey{Namespace: mo.Namespace, Name: mo.Name}, dnSet))
	if err != nil {
		return nil, errors.Wrap(err, "get dn service")
	}

	cnSet := &v1alpha1.CNSet{}
	err, foundCNSet := util.IsFound(ctx.Get(client.ObjectKey{Namespace: mo.Namespace, Name: mo.Name}, cnSet))
	if !foundCNSet && !foundDNSet && !foundLogSet {
		return m.Create, nil
	}

	return nil, nil
}

func (m *MatrixOneActor) Create(ctx *recon.Context[*v1alpha1.MatrixOneCluster]) error {
	mo := ctx.Obj

	logSet := &v1alpha1.LogSet{}
	dnSet := &v1alpha1.DNSet{}
	cnSet := &v1alpha1.CNSet{}
	if err := buildMatrixOneCluster(mo, logSet, dnSet, cnSet); err != nil {
		return err
	}

	return nil
}
