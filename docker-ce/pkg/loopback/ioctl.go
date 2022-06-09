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

package loopback

import (
	"syscall"
	"unsafe"
)

func ioctlLoopCtlGetFree(fd uintptr) (int, error) {
	index, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, LoopCtlGetFree, 0)
	if err != 0 {
		return 0, err
	}
	return int(index), nil
}

func ioctlLoopSetFd(loopFd, sparseFd uintptr) error {
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, loopFd, LoopSetFd, sparseFd); err != 0 {
		return err
	}
	return nil
}

func ioctlLoopSetStatus64(loopFd uintptr, loopInfo *loopInfo64) error {
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, loopFd, LoopSetStatus64, uintptr(unsafe.Pointer(loopInfo))); err != 0 {
		return err
	}
	return nil
}

func ioctlLoopClrFd(loopFd uintptr) error {
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, loopFd, LoopClrFd, 0); err != 0 {
		return err
	}
	return nil
}

func ioctlLoopGetStatus64(loopFd uintptr) (*loopInfo64, error) {
	loopInfo := &loopInfo64{}

	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, loopFd, LoopGetStatus64, uintptr(unsafe.Pointer(loopInfo))); err != 0 {
		return nil, err
	}
	return loopInfo, nil
}

func ioctlLoopSetCapacity(loopFd uintptr, value int) error {
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, loopFd, LoopSetCapacity, uintptr(value)); err != 0 {
		return err
	}
	return nil
}
