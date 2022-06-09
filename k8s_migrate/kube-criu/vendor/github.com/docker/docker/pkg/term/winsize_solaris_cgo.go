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

package term

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

/*
#include <unistd.h>
#include <stropts.h>
#include <termios.h>

// Small wrapper to get rid of variadic args of ioctl()
int my_ioctl(int fd, int cmd, struct winsize *ws) {
	return ioctl(fd, cmd, ws);
}
*/
import "C"

// GetWinsize returns the window size based on the specified file descriptor.
func GetWinsize(fd uintptr) (*Winsize, error) {
	ws := &Winsize{}
	ret, err := C.my_ioctl(C.int(fd), C.int(unix.TIOCGWINSZ), (*C.struct_winsize)(unsafe.Pointer(ws)))
	// Skip retval = 0
	if ret == 0 {
		return ws, nil
	}
	return ws, err
}

// SetWinsize tries to set the specified window size for the specified file descriptor.
func SetWinsize(fd uintptr, ws *Winsize) error {
	ret, err := C.my_ioctl(C.int(fd), C.int(unix.TIOCSWINSZ), (*C.struct_winsize)(unsafe.Pointer(ws)))
	// Skip retval = 0
	if ret == 0 {
		return nil
	}
	return err
}
