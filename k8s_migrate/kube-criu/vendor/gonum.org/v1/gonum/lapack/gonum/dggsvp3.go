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
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/lapack"
)

// Dggsvp3 computes orthogonal matrices U, V and Q such that
//
//                  n-k-l  k    l
//  U^T*A*Q =    k [ 0    A12  A13 ] if m-k-l >= 0;
//               l [ 0     0   A23 ]
//           m-k-l [ 0     0    0  ]
//
//                  n-k-l  k    l
//  U^T*A*Q =    k [ 0    A12  A13 ] if m-k-l < 0;
//             m-k [ 0     0   A23 ]
//
//                  n-k-l  k    l
//  V^T*B*Q =    l [ 0     0   B13 ]
//             p-l [ 0     0    0  ]
//
// where the k×k matrix A12 and l×l matrix B13 are non-singular
// upper triangular. A23 is l×l upper triangular if m-k-l >= 0,
// otherwise A23 is (m-k)×l upper trapezoidal.
//
// Dggsvp3 returns k and l, the dimensions of the sub-blocks. k+l
// is the effective numerical rank of the (m+p)×n matrix [ A^T B^T ]^T.
//
// jobU, jobV and jobQ are options for computing the orthogonal matrices. The behavior
// is as follows
//  jobU == lapack.GSVDU        Compute orthogonal matrix U
//  jobU == lapack.GSVDNone     Do not compute orthogonal matrix.
// The behavior is the same for jobV and jobQ with the exception that instead of
// lapack.GSVDU these accept lapack.GSVDV and lapack.GSVDQ respectively.
// The matrices U, V and Q must be m×m, p×p and n×n respectively unless the
// relevant job parameter is lapack.GSVDNone.
//
// tola and tolb are the convergence criteria for the Jacobi-Kogbetliantz
// iteration procedure. Generally, they are the same as used in the preprocessing
// step, for example,
//  tola = max(m, n)*norm(A)*eps,
//  tolb = max(p, n)*norm(B)*eps.
// Where eps is the machine epsilon.
//
// iwork must have length n, work must have length at least max(1, lwork), and
// lwork must be -1 or greater than zero, otherwise Dggsvp3 will panic.
//
// Dggsvp3 is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dggsvp3(jobU, jobV, jobQ lapack.GSVDJob, m, p, n int, a []float64, lda int, b []float64, ldb int, tola, tolb float64, u []float64, ldu int, v []float64, ldv int, q []float64, ldq int, iwork []int, tau, work []float64, lwork int) (k, l int) {
	const forward = true

	checkMatrix(m, n, a, lda)
	checkMatrix(p, n, b, ldb)

	wantu := jobU == lapack.GSVDU
	if !wantu && jobU != lapack.GSVDNone {
		panic(badGSVDJob + "U")
	}
	if jobU != lapack.GSVDNone {
		checkMatrix(m, m, u, ldu)
	}

	wantv := jobV == lapack.GSVDV
	if !wantv && jobV != lapack.GSVDNone {
		panic(badGSVDJob + "V")
	}
	if jobV != lapack.GSVDNone {
		checkMatrix(p, p, v, ldv)
	}

	wantq := jobQ == lapack.GSVDQ
	if !wantq && jobQ != lapack.GSVDNone {
		panic(badGSVDJob + "Q")
	}
	if jobQ != lapack.GSVDNone {
		checkMatrix(n, n, q, ldq)
	}

	if len(iwork) != n {
		panic(badWork)
	}
	if lwork != -1 && lwork < 1 {
		panic(badWork)
	}
	if len(work) < max(1, lwork) {
		panic(badWork)
	}

	var lwkopt int
	impl.Dgeqp3(p, n, b, ldb, iwork, tau, work, -1)
	lwkopt = int(work[0])
	if wantv {
		lwkopt = max(lwkopt, p)
	}
	lwkopt = max(lwkopt, min(n, p))
	lwkopt = max(lwkopt, m)
	if wantq {
		lwkopt = max(lwkopt, n)
	}
	impl.Dgeqp3(m, n, a, lda, iwork, tau, work, -1)
	lwkopt = max(lwkopt, int(work[0]))
	lwkopt = max(1, lwkopt)
	if lwork == -1 {
		work[0] = float64(lwkopt)
		return 0, 0
	}

	// tau check must come after lwkopt query since
	// the Dggsvd3 call for lwkopt query may have
	// lwork == -1, and tau is provided by work.
	if len(tau) < n {
		panic(badTau)
	}

	// QR with column pivoting of B: B*P = V*[ S11 S12 ].
	//                                       [  0   0  ]
	for i := range iwork[:n] {
		iwork[i] = 0
	}
	impl.Dgeqp3(p, n, b, ldb, iwork, tau, work, lwork)

	// Update A := A*P.
	impl.Dlapmt(forward, m, n, a, lda, iwork)

	// Determine the effective rank of matrix B.
	for i := 0; i < min(p, n); i++ {
		if math.Abs(b[i*ldb+i]) > tolb {
			l++
		}
	}

	if wantv {
		// Copy the details of V, and form V.
		impl.Dlaset(blas.All, p, p, 0, 0, v, ldv)
		if p > 1 {
			impl.Dlacpy(blas.Lower, p-1, min(p, n), b[ldb:], ldb, v[ldv:], ldv)
		}
		impl.Dorg2r(p, p, min(p, n), v, ldv, tau, work)
	}

	// Clean up B.
	for i := 1; i < l; i++ {
		r := b[i*ldb : i*ldb+i]
		for j := range r {
			r[j] = 0
		}
	}
	if p > l {
		impl.Dlaset(blas.All, p-l, n, 0, 0, b[l*ldb:], ldb)
	}

	if wantq {
		// Set Q = I and update Q := Q*P.
		impl.Dlaset(blas.All, n, n, 0, 1, q, ldq)
		impl.Dlapmt(forward, n, n, q, ldq, iwork)
	}

	if p >= l && n != l {
		// RQ factorization of [ S11 S12 ]: [ S11 S12 ] = [ 0 S12 ]*Z.
		impl.Dgerq2(l, n, b, ldb, tau, work)

		// Update A := A*Z^T.
		impl.Dormr2(blas.Right, blas.Trans, m, n, l, b, ldb, tau, a, lda, work)

		if wantq {
			// Update Q := Q*Z^T.
			impl.Dormr2(blas.Right, blas.Trans, n, n, l, b, ldb, tau, q, ldq, work)
		}

		// Clean up B.
		impl.Dlaset(blas.All, l, n-l, 0, 0, b, ldb)
		for i := 1; i < l; i++ {
			r := b[i*ldb+n-l : i*ldb+i+n-l]
			for j := range r {
				r[j] = 0
			}
		}
	}

	// Let              N-L     L
	//            A = [ A11    A12 ] M,
	//
	// then the following does the complete QR decomposition of A11:
	//
	//          A11 = U*[  0  T12 ]*P1^T.
	//                  [  0   0  ]
	for i := range iwork[:n-l] {
		iwork[i] = 0
	}
	impl.Dgeqp3(m, n-l, a, lda, iwork[:n-l], tau, work, lwork)

	// Determine the effective rank of A11.
	for i := 0; i < min(m, n-l); i++ {
		if math.Abs(a[i*lda+i]) > tola {
			k++
		}
	}

	// Update A12 := U^T*A12, where A12 = A[0:m, n-l:n].
	impl.Dorm2r(blas.Left, blas.Trans, m, l, min(m, n-l), a, lda, tau, a[n-l:], lda, work)

	if wantu {
		// Copy the details of U, and form U.
		impl.Dlaset(blas.All, m, m, 0, 0, u, ldu)
		if m > 1 {
			impl.Dlacpy(blas.Lower, m-1, min(m, n-l), a[lda:], lda, u[ldu:], ldu)
		}
		impl.Dorg2r(m, m, min(m, n-l), u, ldu, tau, work)
	}

	if wantq {
		// Update Q[0:n, 0:n-l] := Q[0:n, 0:n-l]*P1.
		impl.Dlapmt(forward, n, n-l, q, ldq, iwork[:n-l])
	}

	// Clean up A: set the strictly lower triangular part of
	// A[0:k, 0:k] = 0, and A[k:m, 0:n-l] = 0.
	for i := 1; i < k; i++ {
		r := a[i*lda : i*lda+i]
		for j := range r {
			r[j] = 0
		}
	}
	if m > k {
		impl.Dlaset(blas.All, m-k, n-l, 0, 0, a[k*lda:], lda)
	}

	if n-l > k {
		// RQ factorization of [ T11 T12 ] = [ 0 T12 ]*Z1.
		impl.Dgerq2(k, n-l, a, lda, tau, work)

		if wantq {
			// Update Q[0:n, 0:n-l] := Q[0:n, 0:n-l]*Z1^T.
			impl.Dorm2r(blas.Right, blas.Trans, n, n-l, k, a, lda, tau, q, ldq, work)
		}

		// Clean up A.
		impl.Dlaset(blas.All, k, n-l-k, 0, 0, a, lda)
		for i := 1; i < k; i++ {
			r := a[i*lda+n-k-l : i*lda+i+n-k-l]
			for j := range r {
				a[j] = 0
			}
		}
	}

	if m > k {
		// QR factorization of A[k:m, n-l:n].
		impl.Dgeqr2(m-k, l, a[k*lda+n-l:], lda, tau, work)
		if wantu {
			// Update U[:, k:m) := U[:, k:m]*U1.
			impl.Dorm2r(blas.Right, blas.NoTrans, m, m-k, min(m-k, l), a[k*lda+n-l:], lda, tau, u[k:], ldu, work)
		}

		// Clean up A.
		for i := k + 1; i < m; i++ {
			r := a[i*lda+n-l : i*lda+min(n-l+i-k, n)]
			for j := range r {
				r[j] = 0
			}
		}
	}

	work[0] = float64(lwkopt)
	return k, l
}
