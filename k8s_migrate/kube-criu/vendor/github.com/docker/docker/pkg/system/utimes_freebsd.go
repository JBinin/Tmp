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
package system

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// LUtimesNano is used to change access and modification time of the specified path.
// It's used for symbol link file because unix.UtimesNano doesn't support a NOFOLLOW flag atm.
func LUtimesNano(path string, ts []syscall.Timespec) error {
	var _path *byte
	_path, err := unix.BytePtrFromString(path)
	if err != nil {
		return err
	}

	if _, _, err := unix.Syscall(unix.SYS_LUTIMES, uintptr(unsafe.Pointer(_path)), uintptr(unsafe.Pointer(&ts[0])), 0); err != 0 && err != unix.ENOSYS {
		return err
	}

	return nil
}
