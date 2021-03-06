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
// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import "gonum.org/v1/gonum/blas"

// Dorgtr generates a real orthogonal matrix Q which is defined as the product
// of n-1 elementary reflectors of order n as returned by Dsytrd.
//
// The construction of Q depends on the value of uplo:
//  Q = H_{n-1} * ... * H_1 * H_0  if uplo == blas.Upper
//  Q = H_0 * H_1 * ... * H_{n-1}  if uplo == blas.Lower
// where H_i is constructed from the elementary reflectors as computed by Dsytrd.
// See the documentation for Dsytrd for more information.
//
// tau must have length at least n-1, and Dorgtr will panic otherwise.
//
// work is temporary storage, and lwork specifies the usable memory length. At
// minimum, lwork >= max(1,n-1), and Dorgtr will panic otherwise. The amount of blocking
// is limited by the usable length.
// If lwork == -1, instead of computing Dorgtr the optimal work length is stored
// into work[0].
//
// Dorgtr is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dorgtr(uplo blas.Uplo, n int, a []float64, lda int, tau, work []float64, lwork int) {
	checkMatrix(n, n, a, lda)
	if len(tau) < n-1 {
		panic(badTau)
	}
	if len(work) < lwork {
		panic(badWork)
	}
	if lwork < n-1 && lwork != -1 {
		panic(badWork)
	}
	upper := uplo == blas.Upper
	if !upper && uplo != blas.Lower {
		panic(badUplo)
	}

	if n == 0 {
		work[0] = 1
		return
	}

	var nb int
	if upper {
		nb = impl.Ilaenv(1, "DORGQL", " ", n-1, n-1, n-1, -1)
	} else {
		nb = impl.Ilaenv(1, "DORGQR", " ", n-1, n-1, n-1, -1)
	}
	lworkopt := max(1, n-1) * nb
	if lwork == -1 {
		work[0] = float64(lworkopt)
		return
	}

	if upper {
		// Q was determined by a call to Dsytrd with uplo == blas.Upper.
		// Shift the vectors which define the elementary reflectors one column
		// to the left, and set the last row and column of Q to those of the unit
		// matrix.
		for j := 0; j < n-1; j++ {
			for i := 0; i < j; i++ {
				a[i*lda+j] = a[i*lda+j+1]
			}
			a[(n-1)*lda+j] = 0
		}
		for i := 0; i < n-1; i++ {
			a[i*lda+n-1] = 0
		}
		a[(n-1)*lda+n-1] = 1

		// Generate Q[0:n-1, 0:n-1].
		impl.Dorgql(n-1, n-1, n-1, a, lda, tau, work, lwork)
	} else {
		// Q was determined by a call to Dsytrd with uplo == blas.Upper.
		// Shift the vectors which define the elementary reflectors one column
		// to the right, and set the first row and column of Q to those of the unit
		// matrix.
		for j := n - 1; j > 0; j-- {
			a[j] = 0
			for i := j + 1; i < n; i++ {
				a[i*lda+j] = a[i*lda+j-1]
			}
		}
		a[0] = 1
		for i := 1; i < n; i++ {
			a[i*lda] = 0
		}
		if n > 1 {
			// Generate Q[1:n, 1:n].
			impl.Dorgqr(n-1, n-1, n-1, a[lda+1:], lda, tau, work, lwork)
		}
	}
	work[0] = float64(lworkopt)
}
