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

// Iladlr scans a matrix for its last non-zero row. Returns -1 if the matrix
// is all zeros.
//
// Iladlr is an internal routine. It is exported for testing purposes.
func (Implementation) Iladlr(m, n int, a []float64, lda int) int {
	if m == 0 {
		return m - 1
	}

	checkMatrix(m, n, a, lda)

	// Check the common case where the corner is non-zero
	if a[(m-1)*lda] != 0 || a[(m-1)*lda+n-1] != 0 {
		return m - 1
	}
	for i := m - 1; i >= 0; i-- {
		for j := 0; j < n; j++ {
			if a[i*lda+j] != 0 {
				return i
			}
		}
	}
	return -1
}
