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
// +build unsafe

// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"unsafe"
)

// This file has unsafe variants of some helper methods.

type unsafeString struct {
	Data uintptr
	Len  int
}

type unsafeSlice struct {
	Data uintptr
	Len  int
	Cap  int
}

// stringView returns a view of the []byte as a string.
// In unsafe mode, it doesn't incur allocation and copying caused by conversion.
// In regular safe mode, it is an allocation and copy.
func stringView(v []byte) string {
	if len(v) == 0 {
		return ""
	}

	bx := (*unsafeSlice)(unsafe.Pointer(&v))
	sx := unsafeString{bx.Data, bx.Len}
	return *(*string)(unsafe.Pointer(&sx))
}

// bytesView returns a view of the string as a []byte.
// In unsafe mode, it doesn't incur allocation and copying caused by conversion.
// In regular safe mode, it is an allocation and copy.
func bytesView(v string) []byte {
	if len(v) == 0 {
		return zeroByteSlice
	}

	sx := (*unsafeString)(unsafe.Pointer(&v))
	bx := unsafeSlice{sx.Data, sx.Len, sx.Len}
	return *(*[]byte)(unsafe.Pointer(&bx))
}
