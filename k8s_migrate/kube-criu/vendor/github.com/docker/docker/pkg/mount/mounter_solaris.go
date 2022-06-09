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
// +build solaris,cgo

package mount

import (
	"golang.org/x/sys/unix"
	"unsafe"
)

// #include <stdlib.h>
// #include <stdio.h>
// #include <sys/mount.h>
// int Mount(const char *spec, const char *dir, int mflag,
// char *fstype, char *dataptr, int datalen, char *optptr, int optlen) {
//     return mount(spec, dir, mflag, fstype, dataptr, datalen, optptr, optlen);
// }
import "C"

func mount(device, target, mType string, flag uintptr, data string) error {
	spec := C.CString(device)
	dir := C.CString(target)
	fstype := C.CString(mType)
	_, err := C.Mount(spec, dir, C.int(flag), fstype, nil, 0, nil, 0)
	C.free(unsafe.Pointer(spec))
	C.free(unsafe.Pointer(dir))
	C.free(unsafe.Pointer(fstype))
	return err
}

func unmount(target string, flag int) error {
	err := unix.Unmount(target, flag)
	return err
}
