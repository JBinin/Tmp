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
// +build linux,386

package system

import (
	"syscall"
)

// Setuid sets the uid of the calling thread to the specified uid.
func Setuid(uid int) (err error) {
	_, _, e1 := syscall.RawSyscall(syscall.SYS_SETUID32, uintptr(uid), 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}

// Setgid sets the gid of the calling thread to the specified gid.
func Setgid(gid int) (err error) {
	_, _, e1 := syscall.RawSyscall(syscall.SYS_SETGID32, uintptr(gid), 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}
