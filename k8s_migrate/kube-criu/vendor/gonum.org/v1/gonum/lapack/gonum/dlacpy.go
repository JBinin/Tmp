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

// Dlacpy copies the elements of A specified by uplo into B. Uplo can specify
// a triangular portion with blas.Upper or blas.Lower, or can specify all of the
// elemest with blas.All.
//
// Dlacpy is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dlacpy(uplo blas.Uplo, m, n int, a []float64, lda int, b []float64, ldb int) {
	checkMatrix(m, n, a, lda)
	checkMatrix(m, n, b, ldb)
	switch uplo {
	default:
		panic(badUplo)
	case blas.Upper:
		for i := 0; i < m; i++ {
			for j := i; j < n; j++ {
				b[i*ldb+j] = a[i*lda+j]
			}
		}

	case blas.Lower:
		for i := 0; i < m; i++ {
			for j := 0; j < min(i+1, n); j++ {
				b[i*ldb+j] = a[i*lda+j]
			}
		}
	case blas.All:
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				b[i*ldb+j] = a[i*lda+j]
			}
		}
	}
}
