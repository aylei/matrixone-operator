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
	"github.com/matrixorigin/matrixone-operator/api/core/v1alpha1"
	recon "github.com/matrixorigin/matrixone-operator/runtime/pkg/reconciler"
	"github.com/matrixorigin/matrixone-operator/runtime/pkg/util"
	kruise "github.com/openkruise/kruise-api/apps/v1alpha1"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Controller struct {
	targetNamespacedName types.NamespacedName
	cloneSet             *kruise.CloneSet
	service              *corev1.Service
	hService             *corev1.Service
}

var _ recon.Actor[*v1alpha1.DNSet] = &Controller{}

func (c *Controller) Observe(ctx *recon.Context[*v1alpha1.DNSet]) (recon.Action[*v1alpha1.DNSet], error) {
	cloneSet := c.cloneSet
	err, foundCs := util.IsFound(ctx.Get(c.targetNamespacedName, cloneSet))
	if err != nil {
		return nil, errors.Wrap(err, "get dn service cloneset")
	}

	if !foundCs {
		return c.Create, nil
	}
	return nil, nil
}

func (c *Controller) Finalize(ctx *recon.Context[*v1alpha1.DNSet]) (bool, error) {
	dn := ctx.Obj
	var errs error

	svcExit, err := ctx.Exist(client.ObjectKey{Namespace: dn.Namespace, Name: dn.Name}, c.service)
	err = multierr.Append(errs, err)
	hSvcExit, err := ctx.Exist(client.ObjectKey{Namespace: dn.Namespace, Name: getDNSetHeadlessSvcName(dn)}, c.hService)
	errs = multierr.Append(errs, err)
	dnSetExit, err := ctx.Exist(client.ObjectKey{Namespace: dn.Namespace, Name: getDNSetName(dn)}, c.cloneSet)
	errs = multierr.Append(errs, err)

	res := !hSvcExit && !dnSetExit && !svcExit

	return res, nil
}

func (c *Controller) Create(ctx *recon.Context[*v1alpha1.DNSet]) error {

	return nil
}
