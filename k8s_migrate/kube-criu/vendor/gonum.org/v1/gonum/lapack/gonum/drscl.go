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
	"math"

	"gonum.org/v1/gonum/blas/blas64"
)

// Drscl multiplies the vector x by 1/a being careful to avoid overflow or
// underflow where possible.
//
// Drscl is an internal routine. It is exported for testing purposes.
func (impl Implementation) Drscl(n int, a float64, x []float64, incX int) {
	checkVector(n, x, incX)
	bi := blas64.Implementation()
	cden := a
	cnum := 1.0
	smlnum := dlamchS
	bignum := 1 / smlnum
	for {
		cden1 := cden * smlnum
		cnum1 := cnum / bignum
		var mul float64
		var done bool
		switch {
		case cnum != 0 && math.Abs(cden1) > math.Abs(cnum):
			mul = smlnum
			done = false
			cden = cden1
		case math.Abs(cnum1) > math.Abs(cden):
			mul = bignum
			done = false
			cnum = cnum1
		default:
			mul = cnum / cden
			done = true
		}
		bi.Dscal(n, mul, x, incX)
		if done {
			break
		}
	}
}
