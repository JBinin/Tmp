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
package pty

import (
	"os"
	"syscall"
	"unsafe"
)

// Getsize returns the number of rows (lines) and cols (positions
// in each line) in terminal t.
func Getsize(t *os.File) (rows, cols int, err error) {
	var ws winsize
	err = windowrect(&ws, t.Fd())
	return int(ws.ws_row), int(ws.ws_col), err
}

type winsize struct {
	ws_row    uint16
	ws_col    uint16
	ws_xpixel uint16
	ws_ypixel uint16
}

func windowrect(ws *winsize, fd uintptr) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(ws)),
	)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}
