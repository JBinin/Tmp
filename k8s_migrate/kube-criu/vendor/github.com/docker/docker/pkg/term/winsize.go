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
// +build !solaris,!windows

package term

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// GetWinsize returns the window size based on the specified file descriptor.
func GetWinsize(fd uintptr) (*Winsize, error) {
	ws := &Winsize{}
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, uintptr(unix.TIOCGWINSZ), uintptr(unsafe.Pointer(ws)))
	// Skipp errno = 0
	if err == 0 {
		return ws, nil
	}
	return ws, err
}

// SetWinsize tries to set the specified window size for the specified file descriptor.
func SetWinsize(fd uintptr, ws *Winsize) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, uintptr(unix.TIOCSWINSZ), uintptr(unsafe.Pointer(ws)))
	// Skipp errno = 0
	if err == 0 {
		return nil
	}
	return err
}
