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

import "github.com/matrixorigin/matrixone-operator/api/core/v1alpha1"

func buildMatrixOneCluster(
	mo *v1alpha1.MatrixOneCluster,
	logSet *v1alpha1.LogSet,
	dn *v1alpha1.DNSet,
	cn *v1alpha1.CNSet) error {
	if err := syncCNBasic(mo, cn); err != nil {
		return err
	}

	if err := syncDNBasic(mo, dn); err != nil {
		return err
	}

	if err := syncLogServiceBasic(mo, logSet); err != nil {
		return err
	}
	return nil
}
