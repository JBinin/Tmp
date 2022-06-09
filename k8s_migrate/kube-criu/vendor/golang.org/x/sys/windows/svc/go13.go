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
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows
// +build go1.3

package svc

import "unsafe"

const ptrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const

// Should be a built-in for unsafe.Pointer?
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

// funcPC returns the entry PC of the function f.
// It assumes that f is a func value. Otherwise the behavior is undefined.
func funcPC(f interface{}) uintptr {
	return **(**uintptr)(add(unsafe.Pointer(&f), ptrSize))
}

// from sys_386.s and sys_amd64.s
func servicectlhandler(ctl uint32) uintptr
func servicemain(argc uint32, argv **uint16)

func getServiceMain(r *uintptr) {
	*r = funcPC(servicemain)
}
