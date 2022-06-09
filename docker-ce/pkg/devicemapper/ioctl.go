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
// +build linux

package devicemapper

import (
	"syscall"
	"unsafe"
)

func ioctlBlkGetSize64(fd uintptr) (int64, error) {
	var size int64
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, BlkGetSize64, uintptr(unsafe.Pointer(&size))); err != 0 {
		return 0, err
	}
	return size, nil
}

func ioctlBlkDiscard(fd uintptr, offset, length uint64) error {
	var r [2]uint64
	r[0] = offset
	r[1] = length

	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, BlkDiscard, uintptr(unsafe.Pointer(&r[0]))); err != 0 {
		return err
	}
	return nil
}
