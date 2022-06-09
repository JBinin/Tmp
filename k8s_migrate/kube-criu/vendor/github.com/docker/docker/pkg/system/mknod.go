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

package system

import (
	"golang.org/x/sys/unix"
)

// Mknod creates a filesystem node (file, device special file or named pipe) named path
// with attributes specified by mode and dev.
func Mknod(path string, mode uint32, dev int) error {
	return unix.Mknod(path, mode, dev)
}

// Mkdev is used to build the value of linux devices (in /dev/) which specifies major
// and minor number of the newly created device special file.
// Linux device nodes are a bit weird due to backwards compat with 16 bit device nodes.
// They are, from low to high: the lower 8 bits of the minor, then 12 bits of the major,
// then the top 12 bits of the minor.
func Mkdev(major int64, minor int64) uint32 {
	return uint32(((minor & 0xfff00) << 12) | ((major & 0xfff) << 8) | (minor & 0xff))
}
