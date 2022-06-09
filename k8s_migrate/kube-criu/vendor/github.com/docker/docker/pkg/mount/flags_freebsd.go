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
// +build freebsd,cgo

package mount

/*
#include <sys/mount.h>
*/
import "C"

const (
	// RDONLY will mount the filesystem as read-only.
	RDONLY = C.MNT_RDONLY

	// NOSUID will not allow set-user-identifier or set-group-identifier bits to
	// take effect.
	NOSUID = C.MNT_NOSUID

	// NOEXEC will not allow execution of any binaries on the mounted file system.
	NOEXEC = C.MNT_NOEXEC

	// SYNCHRONOUS will allow any I/O to the file system to be done synchronously.
	SYNCHRONOUS = C.MNT_SYNCHRONOUS

	// NOATIME will not update the file access time when reading from a file.
	NOATIME = C.MNT_NOATIME
)

// These flags are unsupported.
const (
	BIND        = 0
	DIRSYNC     = 0
	MANDLOCK    = 0
	NODEV       = 0
	NODIRATIME  = 0
	UNBINDABLE  = 0
	RUNBINDABLE = 0
	PRIVATE     = 0
	RPRIVATE    = 0
	SHARED      = 0
	RSHARED     = 0
	SLAVE       = 0
	RSLAVE      = 0
	RBIND       = 0
	RELATIVE    = 0
	RELATIME    = 0
	REMOUNT     = 0
	STRICTATIME = 0
	mntDetach   = 0
)
