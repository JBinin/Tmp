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
// +build !windows
// +build !solaris !cgo

package term

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func tcget(fd uintptr, p *Termios) syscall.Errno {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, uintptr(getTermios), uintptr(unsafe.Pointer(p)))
	return err
}

func tcset(fd uintptr, p *Termios) syscall.Errno {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, setTermios, uintptr(unsafe.Pointer(p)))
	return err
}
