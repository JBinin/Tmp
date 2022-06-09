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

package c128

// ScalUnitaryTo is
//  for i, v := range x {
//  	dst[i] = alpha * v
//  }
func ScalUnitaryTo(dst []complex128, alpha complex128, x []complex128) {
	for i, v := range x {
		dst[i] = alpha * v
	}
}

// ScalIncTo is
//  var idst, ix uintptr
//  for i := 0; i < int(n); i++ {
//  	dst[idst] = alpha * x[ix]
//  	ix += incX
//  	idst += incDst
//  }
func ScalIncTo(dst []complex128, incDst uintptr, alpha complex128, x []complex128, n, incX uintptr) {
	var idst, ix uintptr
	for i := 0; i < int(n); i++ {
		dst[idst] = alpha * x[ix]
		ix += incX
		idst += incDst
	}
}
