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

// +build !noasm,!appengine,!safe

package f64

// Ger performs the rank-one operation
//  A += alpha * x * y^T
// where A is an m×n dense matrix, x and y are vectors, and alpha is a scalar.
func Ger(m, n uintptr, alpha float64, x []float64, incX uintptr, y []float64, incY uintptr, a []float64, lda uintptr)

// GemvN computes
//  y = alpha * A * x + beta * y
// where A is an m×n dense matrix, x and y are vectors, and alpha and beta are scalars.
func GemvN(m, n uintptr, alpha float64, a []float64, lda uintptr, x []float64, incX uintptr, beta float64, y []float64, incY uintptr)

// GemvT computes
//  y = alpha * A^T * x + beta * y
// where A is an m×n dense matrix, x and y are vectors, and alpha and beta are scalars.
func GemvT(m, n uintptr, alpha float64, a []float64, lda uintptr, x []float64, incX uintptr, beta float64, y []float64, incY uintptr)
