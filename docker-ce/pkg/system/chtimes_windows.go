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
// +build windows

package system

import (
	"syscall"
	"time"
)

//setCTime will set the create time on a file. On Windows, this requires
//calling SetFileTime and explicitly including the create time.
func setCTime(path string, ctime time.Time) error {
	ctimespec := syscall.NsecToTimespec(ctime.UnixNano())
	pathp, e := syscall.UTF16PtrFromString(path)
	if e != nil {
		return e
	}
	h, e := syscall.CreateFile(pathp,
		syscall.FILE_WRITE_ATTRIBUTES, syscall.FILE_SHARE_WRITE, nil,
		syscall.OPEN_EXISTING, syscall.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if e != nil {
		return e
	}
	defer syscall.Close(h)
	c := syscall.NsecToFiletime(syscall.TimespecToNsec(ctimespec))
	return syscall.SetFileTime(h, &c, nil, nil)
}
