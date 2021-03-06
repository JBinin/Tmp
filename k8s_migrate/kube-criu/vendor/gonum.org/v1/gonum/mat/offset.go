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

// +build !appengine,!safe

package mat

import "unsafe"

// offset returns the number of float64 values b[0] is after a[0].
func offset(a, b []float64) int {
	if &a[0] == &b[0] {
		return 0
	}
	// This expression must be atomic with respect to GC moves.
	// At this stage this is true, because the GC does not
	// move. See https://golang.org/issue/12445.
	return int(uintptr(unsafe.Pointer(&b[0]))-uintptr(unsafe.Pointer(&a[0]))) / int(unsafe.Sizeof(float64(0)))
}
