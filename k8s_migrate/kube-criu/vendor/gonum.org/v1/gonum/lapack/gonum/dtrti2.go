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

// Dtrti2 computes the inverse of a triangular matrix, storing the result in place
// into a. This is the BLAS level 2 version of the algorithm.
//
// Dtrti2 is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dtrti2(uplo blas.Uplo, diag blas.Diag, n int, a []float64, lda int) {
	checkMatrix(n, n, a, lda)
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	if diag != blas.NonUnit && diag != blas.Unit {
		panic(badDiag)
	}
	bi := blas64.Implementation()

	nonUnit := diag == blas.NonUnit
	// TODO(btracey): Replace this with a row-major ordering.
	if uplo == blas.Upper {
		for j := 0; j < n; j++ {
			var ajj float64
			if nonUnit {
				ajj = 1 / a[j*lda+j]
				a[j*lda+j] = ajj
				ajj *= -1
			} else {
				ajj = -1
			}
			bi.Dtrmv(blas.Upper, blas.NoTrans, diag, j, a, lda, a[j:], lda)
			bi.Dscal(j, ajj, a[j:], lda)
		}
		return
	}
	for j := n - 1; j >= 0; j-- {
		var ajj float64
		if nonUnit {
			ajj = 1 / a[j*lda+j]
			a[j*lda+j] = ajj
			ajj *= -1
		} else {
			ajj = -1
		}
		if j < n-1 {
			bi.Dtrmv(blas.Lower, blas.NoTrans, diag, n-j-1, a[(j+1)*lda+j+1:], lda, a[(j+1)*lda+j:], lda)
			bi.Dscal(n-j-1, ajj, a[(j+1)*lda+j:], lda)
		}
	}
}
