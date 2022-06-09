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
// +build !linux,!freebsd freebsd,!cgo solaris,!cgo

package mount

// These flags are unsupported.
const (
	BIND        = 0
	DIRSYNC     = 0
	MANDLOCK    = 0
	NOATIME     = 0
	NODEV       = 0
	NODIRATIME  = 0
	NOEXEC      = 0
	NOSUID      = 0
	UNBINDABLE  = 0
	RUNBINDABLE = 0
	PRIVATE     = 0
	RPRIVATE    = 0
	SHARED      = 0
	RSHARED     = 0
	SLAVE       = 0
	RSLAVE      = 0
	RBIND       = 0
	RELATIME    = 0
	RELATIVE    = 0
	REMOUNT     = 0
	STRICTATIME = 0
	SYNCHRONOUS = 0
	RDONLY      = 0
)
