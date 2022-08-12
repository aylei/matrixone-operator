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
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Actor struct {
	Mgr    manager.Manager
	LActor *logset.LogSetActor
	DActor *dnset.DNSetActor
	CActor *cnset.CNSetActor
}

var _ recon.Actor[*v1alpha1.MatrixOneCluster] = &Actor{}

func (m *Actor) Finalize(ctx *recon.Context[*v1alpha1.MatrixOneCluster]) (bool, error) {
	mo := ctx.Obj

	klog.V(recon.Info).Info("finalize matrixone")

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

func (m *Actor) Observe(
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
	if err != nil {
		return nil, errors.Wrap(err, "get cn service")
	}

	if !foundCNSet && !foundDNSet && !foundLogSet {
		return m.Create, nil
	}

	return nil, nil
}

func (m *Actor) Create(ctx *recon.Context[*v1alpha1.MatrixOneCluster]) error {
	mo := ctx.Obj
	logSet := &v1alpha1.LogSet{}
	dnSet := &v1alpha1.DNSet{}
	cnSet := &v1alpha1.CNSet{}

	klog.V(recon.Info).Info("create matrixone cluster...")

	if err := buildMatrixOneCluster(mo, logSet, dnSet, cnSet); err != nil {
		return err
	}

	return nil
}

func (m *Actor) Reconcile() error {
	logSet := &v1alpha1.LogSet{}
	dnSet := &v1alpha1.DNSet{}
	cnSet := &v1alpha1.CNSet{}
	ms := &v1alpha1.MatrixOneCluster{}

	if err := recon.Setup[*v1alpha1.MatrixOneCluster](ms, "matrixone", m.Mgr, m); err != nil {
		return err
	}

	if err := initialDNSet(m.Mgr, dnSet, m.DActor); err != nil {
		return err
	}
	if err := initialCNSet(m.Mgr, cnSet, m.CActor); err != nil {
		return err
	}
	if err := initialLogSet(m.Mgr, logSet, m.LActor); err != nil {
		return err
	}

	return nil
}

func initialDNSet(mgr manager.Manager, dn *v1alpha1.DNSet, actor *dnset.DNSetActor) error {
	if err := actor.Reconcile(mgr, dn); err != nil {
		return err
	}

	return nil
}

func initialCNSet(mgr manager.Manager, cn *v1alpha1.CNSet, actor *cnset.CNSetActor) error {
	if err := actor.Reconcile(mgr, cn); err != nil {
		return err
	}

	return nil
}

func initialLogSet(mgr manager.Manager, ls *v1alpha1.LogSet, actor *logset.LogSetActor) error {
	if err := actor.Reconcile(mgr, ls); err != nil {
		return err
	}

	return nil
}
