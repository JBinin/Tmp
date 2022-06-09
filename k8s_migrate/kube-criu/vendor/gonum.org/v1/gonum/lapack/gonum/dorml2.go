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

import "gonum.org/v1/gonum/blas"

// Dorml2 multiplies a general matrix C by an orthogonal matrix from an LQ factorization
// determined by Dgelqf.
//  C = Q * C    if side == blas.Left and trans == blas.NoTrans
//  C = Q^T * C  if side == blas.Left and trans == blas.Trans
//  C = C * Q    if side == blas.Right and trans == blas.NoTrans
//  C = C * Q^T  if side == blas.Right and trans == blas.Trans
// If side == blas.Left, a is a matrix of side k×m, and if side == blas.Right
// a is of size k×n.
//
// tau contains the Householder factors and is of length at least k and this function will
// panic otherwise.
//
// work is temporary storage of length at least n if side == blas.Left
// and at least m if side == blas.Right and this function will panic otherwise.
//
// Dorml2 is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dorml2(side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64) {
	if side != blas.Left && side != blas.Right {
		panic(badSide)
	}
	if trans != blas.Trans && trans != blas.NoTrans {
		panic(badTrans)
	}

	left := side == blas.Left
	notran := trans == blas.NoTrans
	if left {
		checkMatrix(k, m, a, lda)
		if len(work) < n {
			panic(badWork)
		}
	} else {
		checkMatrix(k, n, a, lda)
		if len(work) < m {
			panic(badWork)
		}
	}
	checkMatrix(m, n, c, ldc)
	if m == 0 || n == 0 || k == 0 {
		return
	}
	switch {
	case left && notran:
		for i := 0; i < k; i++ {
			aii := a[i*lda+i]
			a[i*lda+i] = 1
			impl.Dlarf(side, m-i, n, a[i*lda+i:], 1, tau[i], c[i*ldc:], ldc, work)
			a[i*lda+i] = aii
		}

	case left && !notran:
		for i := k - 1; i >= 0; i-- {
			aii := a[i*lda+i]
			a[i*lda+i] = 1
			impl.Dlarf(side, m-i, n, a[i*lda+i:], 1, tau[i], c[i*ldc:], ldc, work)
			a[i*lda+i] = aii
		}

	case !left && notran:
		for i := k - 1; i >= 0; i-- {
			aii := a[i*lda+i]
			a[i*lda+i] = 1
			impl.Dlarf(side, m, n-i, a[i*lda+i:], 1, tau[i], c[i:], ldc, work)
			a[i*lda+i] = aii
		}

	case !left && !notran:
		for i := 0; i < k; i++ {
			aii := a[i*lda+i]
			a[i*lda+i] = 1
			impl.Dlarf(side, m, n-i, a[i*lda+i:], 1, tau[i], c[i:], ldc, work)
			a[i*lda+i] = aii
		}
	}
}
