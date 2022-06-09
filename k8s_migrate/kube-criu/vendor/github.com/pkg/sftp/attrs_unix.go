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
// +build darwin dragonfly freebsd !android,linux netbsd openbsd solaris
// +build cgo

package sftp

import (
	"os"
	"syscall"
)

func fileStatFromInfoOs(fi os.FileInfo, flags *uint32, fileStat *FileStat) {
	if statt, ok := fi.Sys().(*syscall.Stat_t); ok {
		*flags |= ssh_FILEXFER_ATTR_UIDGID
		fileStat.UID = statt.Uid
		fileStat.GID = statt.Gid
	}
}
