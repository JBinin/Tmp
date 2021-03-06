/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

// Dtrtrs solves a triangular system of the form A * X = B or A^T * X = B. Dtrtrs
// returns whether the solve completed successfully. If A is singular, no solve is performed.
func (impl Implementation) Dtrtrs(uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, n, nrhs int, a []float64, lda int, b []float64, ldb int) (ok bool) {
	nounit := diag == blas.NonUnit
	if n == 0 {
		return false
	}
	// Check for singularity.
	if nounit {
		for i := 0; i < n; i++ {
			if a[i*lda+i] == 0 {
				return false
			}
		}
	}
	bi := blas64.Implementation()
	bi.Dtrsm(blas.Left, uplo, trans, diag, n, nrhs, 1, a, lda, b, ldb)
	return true
}
