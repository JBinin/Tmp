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

// +build !amd64 noasm appengine safe

package f64

import "math"

// L1Norm is
//  for _, v := range x {
//  	sum += math.Abs(v)
//  }
//  return sum
func L1Norm(x []float64) (sum float64) {
	for _, v := range x {
		sum += math.Abs(v)
	}
	return sum
}

// L1NormInc is
//  for i := 0; i < n*incX; i += incX {
//  	sum += math.Abs(x[i])
//  }
//  return sum
func L1NormInc(x []float64, n, incX int) (sum float64) {
	for i := 0; i < n*incX; i += incX {
		sum += math.Abs(x[i])
	}
	return sum
}

// Add is
//  for i, v := range s {
//  	dst[i] += v
//  }
func Add(dst, s []float64) {
	for i, v := range s {
		dst[i] += v
	}
}

// AddConst is
//  for i := range x {
//  	x[i] += alpha
//  }
func AddConst(alpha float64, x []float64) {
	for i := range x {
		x[i] += alpha
	}
}

// CumSum is
//  if len(s) == 0 {
//  	return dst
//  }
//  dst[0] = s[0]
//  for i, v := range s[1:] {
//  	dst[i+1] = dst[i] + v
//  }
//  return dst
func CumSum(dst, s []float64) []float64 {
	if len(s) == 0 {
		return dst
	}
	dst[0] = s[0]
	for i, v := range s[1:] {
		dst[i+1] = dst[i] + v
	}
	return dst
}

// CumProd is
//  if len(s) == 0 {
//  	return dst
//  }
//  dst[0] = s[0]
//  for i, v := range s[1:] {
//  	dst[i+1] = dst[i] * v
//  }
//  return dst
func CumProd(dst, s []float64) []float64 {
	if len(s) == 0 {
		return dst
	}
	dst[0] = s[0]
	for i, v := range s[1:] {
		dst[i+1] = dst[i] * v
	}
	return dst
}

// Div is
//  for i, v := range s {
//  	dst[i] /= v
//  }
func Div(dst, s []float64) {
	for i, v := range s {
		dst[i] /= v
	}
}

// DivTo is
//  for i, v := range s {
//  	dst[i] = v / t[i]
//  }
//  return dst
func DivTo(dst, s, t []float64) []float64 {
	for i, v := range s {
		dst[i] = v / t[i]
	}
	return dst
}

// L1Dist is
//  var norm float64
//  for i, v := range s {
//  	norm += math.Abs(t[i] - v)
//  }
//  return norm
func L1Dist(s, t []float64) float64 {
	var norm float64
	for i, v := range s {
		norm += math.Abs(t[i] - v)
	}
	return norm
}

// LinfDist is
//  var norm float64
//  if len(s) == 0 {
//  	return 0
//  }
//  norm = math.Abs(t[0] - s[0])
//  for i, v := range s[1:] {
//  	absDiff := math.Abs(t[i+1] - v)
//  	if absDiff > norm || math.IsNaN(norm) {
//  		norm = absDiff
//  	}
//  }
//  return norm
func LinfDist(s, t []float64) float64 {
	var norm float64
	if len(s) == 0 {
		return 0
	}
	norm = math.Abs(t[0] - s[0])
	for i, v := range s[1:] {
		absDiff := math.Abs(t[i+1] - v)
		if absDiff > norm || math.IsNaN(norm) {
			norm = absDiff
		}
	}
	return norm
}
