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
// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/lapack"
)

// Dgerqf computes an RQ factorization of the m×n matrix A,
//  A = R * Q.
// On exit, if m <= n, the upper triangle of the subarray
// A[0:m, n-m:n] contains the m×m upper triangular matrix R.
// If m >= n, the elements on and above the (m-n)-th subdiagonal
// contain the m×n upper trapezoidal matrix R.
// The remaining elements, with tau, represent the
// orthogonal matrix Q as a product of min(m,n) elementary
// reflectors.
//
// The matrix Q is represented as a product of elementary reflectors
//  Q = H_0 H_1 . . . H_{min(m,n)-1}.
// Each H(i) has the form
//  H_i = I - tau_i * v * v^T
// where v is a vector with v[0:n-k+i-1] stored in A[m-k+i, 0:n-k+i-1],
// v[n-k+i:n] = 0 and v[n-k+i] = 1.
//
// tau must have length min(m,n), work must have length max(1, lwork),
// and lwork must be -1 or at least max(1, m), otherwise Dgerqf will panic.
// On exit, work[0] will contain the optimal length for work.
//
// Dgerqf is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dgerqf(m, n int, a []float64, lda int, tau, work []float64, lwork int) {
	checkMatrix(m, n, a, lda)

	if len(work) < max(1, lwork) {
		panic(shortWork)
	}
	if lwork != -1 && lwork < max(1, m) {
		panic(badWork)
	}

	k := min(m, n)
	if len(tau) != k {
		panic(badTau)
	}

	var nb, lwkopt int
	if k == 0 {
		lwkopt = 1
	} else {
		nb = impl.Ilaenv(1, "DGERQF", " ", m, n, -1, -1)
		lwkopt = m * nb
	}
	work[0] = float64(lwkopt)

	if lwork == -1 {
		return
	}

	// Return quickly if possible.
	if k == 0 {
		return
	}

	nbmin := 2
	nx := 1
	iws := m
	var ldwork int
	if 1 < nb && nb < k {
		// Determine when to cross over from blocked to unblocked code.
		nx = max(0, impl.Ilaenv(3, "DGERQF", " ", m, n, -1, -1))
		if nx < k {
			// Determine whether workspace is large enough for blocked code.
			iws = m * nb
			if lwork < iws {
				// Not enough workspace to use optimal nb. Reduce
				// nb and determine the minimum value of nb.
				nb = lwork / m
				nbmin = max(2, impl.Ilaenv(2, "DGERQF", " ", m, n, -1, -1))
			}
			ldwork = nb
		}
	}

	var mu, nu int
	if nbmin <= nb && nb < k && nx < k {
		// Use blocked code initially.
		// The last kk rows are handled by the block method.
		ki := ((k - nx - 1) / nb) * nb
		kk := min(k, ki+nb)

		var i int
		for i = k - kk + ki; i >= k-kk; i -= nb {
			ib := min(k-i, nb)

			// Compute the RQ factorization of the current block
			// A[m-k+i:m-k+i+ib-1, 0:n-k+i+ib-1].
			impl.Dgerq2(ib, n-k+i+ib, a[(m-k+i)*lda:], lda, tau[i:], work)
			if m-k+i > 0 {
				// Form the triangular factor of the block reflector
				// H = H_{i+ib-1} . . . H_{i+1} H_i.
				impl.Dlarft(lapack.Backward, lapack.RowWise,
					n-k+i+ib, ib, a[(m-k+i)*lda:], lda, tau[i:],
					work, ldwork)

				// Apply H to A[0:m-k+i-1, 0:n-k+i+ib-1] from the right.
				impl.Dlarfb(blas.Right, blas.NoTrans, lapack.Backward, lapack.RowWise,
					m-k+i, n-k+i+ib, ib, a[(m-k+i)*lda:], lda,
					work, ldwork,
					a, lda,
					work[ib*ldwork:], ldwork)
			}
		}
		mu = m - k + i + nb
		nu = n - k + i + nb
	} else {
		mu = m
		nu = n
	}

	// Use unblocked code to factor the last or only block.
	if mu > 0 && nu > 0 {
		impl.Dgerq2(mu, nu, a, lda, tau, work)
	}
	work[0] = float64(iws)
}
