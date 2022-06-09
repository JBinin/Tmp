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
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build openbsd
// +build 386 amd64 arm

package unix

import (
	"syscall"
	"unsafe"
)

const (
	SYS_PLEDGE = 108
)

// Pledge implements the pledge syscall. For more information see pledge(2).
func Pledge(promises string, paths []string) error {
	promisesPtr, err := syscall.BytePtrFromString(promises)
	if err != nil {
		return err
	}
	promisesUnsafe, pathsUnsafe := unsafe.Pointer(promisesPtr), unsafe.Pointer(nil)
	if paths != nil {
		var pathsPtr []*byte
		if pathsPtr, err = syscall.SlicePtrFromStrings(paths); err != nil {
			return err
		}
		pathsUnsafe = unsafe.Pointer(&pathsPtr[0])
	}
	_, _, e := syscall.Syscall(SYS_PLEDGE, uintptr(promisesUnsafe), uintptr(pathsUnsafe), 0)
	if e != 0 {
		return e
	}
	return nil
}
