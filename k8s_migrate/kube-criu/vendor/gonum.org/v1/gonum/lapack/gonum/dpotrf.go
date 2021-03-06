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

// Dpotrf computes the Cholesky decomposition of the symmetric positive definite
// matrix a. If ul == blas.Upper, then a is stored as an upper-triangular matrix,
// and a = U^T U is stored in place into a. If ul == blas.Lower, then a = L L^T
// is computed and stored in-place into a. If a is not positive definite, false
// is returned. This is the blocked version of the algorithm.
func (impl Implementation) Dpotrf(ul blas.Uplo, n int, a []float64, lda int) (ok bool) {
	if ul != blas.Upper && ul != blas.Lower {
		panic(badUplo)
	}
	checkMatrix(n, n, a, lda)

	if n == 0 {
		return true
	}

	nb := impl.Ilaenv(1, "DPOTRF", string(ul), n, -1, -1, -1)
	if nb <= 1 || n <= nb {
		return impl.Dpotf2(ul, n, a, lda)
	}
	bi := blas64.Implementation()
	if ul == blas.Upper {
		for j := 0; j < n; j += nb {
			jb := min(nb, n-j)
			bi.Dsyrk(blas.Upper, blas.Trans, jb, j,
				-1, a[j:], lda,
				1, a[j*lda+j:], lda)
			ok = impl.Dpotf2(blas.Upper, jb, a[j*lda+j:], lda)
			if !ok {
				return ok
			}
			if j+jb < n {
				bi.Dgemm(blas.Trans, blas.NoTrans, jb, n-j-jb, j,
					-1, a[j:], lda, a[j+jb:], lda,
					1, a[j*lda+j+jb:], lda)
				bi.Dtrsm(blas.Left, blas.Upper, blas.Trans, blas.NonUnit, jb, n-j-jb,
					1, a[j*lda+j:], lda,
					a[j*lda+j+jb:], lda)
			}
		}
		return true
	}
	for j := 0; j < n; j += nb {
		jb := min(nb, n-j)
		bi.Dsyrk(blas.Lower, blas.NoTrans, jb, j,
			-1, a[j*lda:], lda,
			1, a[j*lda+j:], lda)
		ok := impl.Dpotf2(blas.Lower, jb, a[j*lda+j:], lda)
		if !ok {
			return ok
		}
		if j+jb < n {
			bi.Dgemm(blas.NoTrans, blas.Trans, n-j-jb, jb, j,
				-1, a[(j+jb)*lda:], lda, a[j*lda:], lda,
				1, a[(j+jb)*lda+j:], lda)
			bi.Dtrsm(blas.Right, blas.Lower, blas.Trans, blas.NonUnit, n-j-jb, jb,
				1, a[j*lda+j:], lda,
				a[(j+jb)*lda+j:], lda)
		}
	}
	return true
}
