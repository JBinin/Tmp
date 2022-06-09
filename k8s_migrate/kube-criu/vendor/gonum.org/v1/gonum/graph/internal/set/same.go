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
// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine,!safe

package set

import "unsafe"

// same determines whether two sets are backed by the same store. In the
// current implementation using hash maps it makes use of the fact that
// hash maps are passed as a pointer to a runtime Hmap struct. A map is
// not seen by the runtime as a pointer though, so we use unsafe to get
// the maps' pointer values to compare.
func same(a, b Nodes) bool {
	return *(*uintptr)(unsafe.Pointer(&a)) == *(*uintptr)(unsafe.Pointer(&b))
}

// intsSame determines whether two sets are backed by the same store. In the
// current implementation using hash maps it makes use of the fact that
// hash maps are passed as a pointer to a runtime Hmap struct. A map is
// not seen by the runtime as a pointer though, so we use unsafe to get
// the maps' pointer values to compare.
func intsSame(a, b Ints) bool {
	return *(*uintptr)(unsafe.Pointer(&a)) == *(*uintptr)(unsafe.Pointer(&b))
}

// int64sSame determines whether two sets are backed by the same store. In the
// current implementation using hash maps it makes use of the fact that
// hash maps are passed as a pointer to a runtime Hmap struct. A map is
// not seen by the runtime as a pointer though, so we use unsafe to get
// the maps' pointer values to compare.
func int64sSame(a, b Int64s) bool {
	return *(*uintptr)(unsafe.Pointer(&a)) == *(*uintptr)(unsafe.Pointer(&b))
}
