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
)

type MatrixOneActor struct {
	dnSet  dnset.DNSetActor
	logSet logset.LogSetActor
	cnSet  cnset.CNSetActor
}

var _ recon.Actor[*v1alpha1.MatrixOneCluster] = &MatrixOneActor{}

func (m *MatrixOneActor) Finalize(ctx *recon.Context[*v1alpha1.MatrixOneCluster]) (bool, error) {

	return true, nil
}

func (m *MatrixOneActor) Observe(
	ctx *recon.Context[*v1alpha1.MatrixOneCluster]) (recon.Action[*v1alpha1.MatrixOneCluster], error) {
	return nil, nil
}

func (m *MatrixOneActor) Create(ctx *recon.Context[*v1alpha1.MatrixOneCluster]) error {
	return nil
}
