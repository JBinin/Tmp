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
)

// LUtimesNano is used to change access and modification time of the specified path.
// It's used for symbol link file because syscall.UtimesNano doesn't support a NOFOLLOW flag atm.
func LUtimesNano(path string, ts []syscall.Timespec) error {
	// These are not currently available in syscall
	atFdCwd := -100
	atSymLinkNoFollow := 0x100

	var _path *byte
	_path, err := syscall.BytePtrFromString(path)
	if err != nil {
		return err
	}

	if _, _, err := syscall.Syscall6(syscall.SYS_UTIMENSAT, uintptr(atFdCwd), uintptr(unsafe.Pointer(_path)), uintptr(unsafe.Pointer(&ts[0])), uintptr(atSymLinkNoFollow), 0, 0); err != 0 && err != syscall.ENOSYS {
		return err
	}

	return nil
}
