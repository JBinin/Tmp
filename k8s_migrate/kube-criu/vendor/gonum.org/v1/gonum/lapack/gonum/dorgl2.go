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

// Dorgl2 generates an m×n matrix Q with orthonormal rows defined by the
// first m rows product of elementary reflectors as computed by Dgelqf.
//  Q = H_0 * H_1 * ... * H_{k-1}
// len(tau) >= k, 0 <= k <= m, 0 <= m <= n, len(work) >= m.
// Dorgl2 will panic if these conditions are not met.
//
// Dorgl2 is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dorgl2(m, n, k int, a []float64, lda int, tau, work []float64) {
	checkMatrix(m, n, a, lda)
	if len(tau) < k {
		panic(badTau)
	}
	if k > m {
		panic(kGTM)
	}
	if k > m {
		panic(kGTM)
	}
	if m > n {
		panic(nLTM)
	}
	if len(work) < m {
		panic(badWork)
	}
	if m == 0 {
		return
	}
	bi := blas64.Implementation()
	if k < m {
		for i := k; i < m; i++ {
			for j := 0; j < n; j++ {
				a[i*lda+j] = 0
			}
		}
		for j := k; j < m; j++ {
			a[j*lda+j] = 1
		}
	}
	for i := k - 1; i >= 0; i-- {
		if i < n-1 {
			if i < m-1 {
				a[i*lda+i] = 1
				impl.Dlarf(blas.Right, m-i-1, n-i, a[i*lda+i:], 1, tau[i], a[(i+1)*lda+i:], lda, work)
			}
			bi.Dscal(n-i-1, -tau[i], a[i*lda+i+1:], 1)
		}
		a[i*lda+i] = 1 - tau[i]
		for l := 0; l < i; l++ {
			a[i*lda+l] = 0
		}
	}
}
