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
	"gonum.org/v1/gonum/lapack"
)

// Dormlq multiplies the matrix C by the orthogonal matrix Q defined by the
// slices a and tau. A and tau are as returned from Dgelqf.
//  C = Q * C    if side == blas.Left and trans == blas.NoTrans
//  C = Q^T * C  if side == blas.Left and trans == blas.Trans
//  C = C * Q    if side == blas.Right and trans == blas.NoTrans
//  C = C * Q^T  if side == blas.Right and trans == blas.Trans
// If side == blas.Left, A is a matrix of side k×m, and if side == blas.Right
// A is of size k×n. This uses a blocked algorithm.
//
// work is temporary storage, and lwork specifies the usable memory length.
// At minimum, lwork >= m if side == blas.Left and lwork >= n if side == blas.Right,
// and this function will panic otherwise.
// Dormlq uses a block algorithm, but the block size is limited
// by the temporary space available. If lwork == -1, instead of performing Dormlq,
// the optimal work length will be stored into work[0].
//
// tau contains the Householder scales and must have length at least k, and
// this function will panic otherwise.
func (impl Implementation) Dormlq(side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64, lwork int) {
	if side != blas.Left && side != blas.Right {
		panic(badSide)
	}
	if trans != blas.Trans && trans != blas.NoTrans {
		panic(badTrans)
	}
	left := side == blas.Left
	if left {
		checkMatrix(k, m, a, lda)
	} else {
		checkMatrix(k, n, a, lda)
	}
	checkMatrix(m, n, c, ldc)
	if len(tau) < k {
		panic(badTau)
	}
	if len(work) < lwork {
		panic(shortWork)
	}
	nw := m
	if left {
		nw = n
	}
	if lwork < max(1, nw) && lwork != -1 {
		panic(badWork)
	}

	if m == 0 || n == 0 || k == 0 {
		work[0] = 1
		return
	}

	const (
		nbmax = 64
		ldt   = nbmax
		tsize = nbmax * ldt
	)
	opts := string(side) + string(trans)
	nb := min(nbmax, impl.Ilaenv(1, "DORMLQ", opts, m, n, k, -1))
	lworkopt := max(1, nw)*nb + tsize
	if lwork == -1 {
		work[0] = float64(lworkopt)
		return
	}

	nbmin := 2
	if 1 < nb && nb < k {
		iws := nw*nb + tsize
		if lwork < iws {
			nb = (lwork - tsize) / nw
			nbmin = max(2, impl.Ilaenv(2, "DORMLQ", opts, m, n, k, -1))
		}
	}
	if nb < nbmin || k <= nb {
		// Call unblocked code.
		impl.Dorml2(side, trans, m, n, k, a, lda, tau, c, ldc, work)
		work[0] = float64(lworkopt)
		return
	}

	t := work[:tsize]
	wrk := work[tsize:]
	ldwrk := nb

	notran := trans == blas.NoTrans
	transt := blas.NoTrans
	if notran {
		transt = blas.Trans
	}

	switch {
	case left && notran:
		for i := 0; i < k; i += nb {
			ib := min(nb, k-i)
			impl.Dlarft(lapack.Forward, lapack.RowWise, m-i, ib,
				a[i*lda+i:], lda,
				tau[i:],
				t, ldt)
			impl.Dlarfb(side, transt, lapack.Forward, lapack.RowWise, m-i, n, ib,
				a[i*lda+i:], lda,
				t, ldt,
				c[i*ldc:], ldc,
				wrk, ldwrk)
		}

	case left && !notran:
		for i := ((k - 1) / nb) * nb; i >= 0; i -= nb {
			ib := min(nb, k-i)
			impl.Dlarft(lapack.Forward, lapack.RowWise, m-i, ib,
				a[i*lda+i:], lda,
				tau[i:],
				t, ldt)
			impl.Dlarfb(side, transt, lapack.Forward, lapack.RowWise, m-i, n, ib,
				a[i*lda+i:], lda,
				t, ldt,
				c[i*ldc:], ldc,
				wrk, ldwrk)
		}

	case !left && notran:
		for i := ((k - 1) / nb) * nb; i >= 0; i -= nb {
			ib := min(nb, k-i)
			impl.Dlarft(lapack.Forward, lapack.RowWise, n-i, ib,
				a[i*lda+i:], lda,
				tau[i:],
				t, ldt)
			impl.Dlarfb(side, transt, lapack.Forward, lapack.RowWise, m, n-i, ib,
				a[i*lda+i:], lda,
				t, ldt,
				c[i:], ldc,
				wrk, ldwrk)
		}

	case !left && !notran:
		for i := 0; i < k; i += nb {
			ib := min(nb, k-i)
			impl.Dlarft(lapack.Forward, lapack.RowWise, n-i, ib,
				a[i*lda+i:], lda,
				tau[i:],
				t, ldt)
			impl.Dlarfb(side, transt, lapack.Forward, lapack.RowWise, m, n-i, ib,
				a[i*lda+i:], lda,
				t, ldt,
				c[i:], ldc,
				wrk, ldwrk)
		}
	}
	work[0] = float64(lworkopt)
}
