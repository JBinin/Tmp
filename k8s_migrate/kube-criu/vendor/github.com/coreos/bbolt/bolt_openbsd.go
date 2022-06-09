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
package bolt

import (
	"syscall"
	"unsafe"
)

const (
	msAsync      = 1 << iota // perform asynchronous writes
	msSync                   // perform synchronous writes
	msInvalidate             // invalidate cached data
)

func msync(db *DB) error {
	_, _, errno := syscall.Syscall(syscall.SYS_MSYNC, uintptr(unsafe.Pointer(db.data)), uintptr(db.datasz), msInvalidate)
	if errno != 0 {
		return errno
	}
	return nil
}

func fdatasync(db *DB) error {
	if db.data != nil {
		return msync(db)
	}
	return db.file.Sync()
}
