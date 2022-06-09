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
	"syscall"
)

// StatT type contains status of a file. It contains metadata
// like permission, owner, group, size, etc about a file.
type StatT struct {
	mode uint32
	uid  uint32
	gid  uint32
	rdev uint64
	size int64
	mtim syscall.Timespec
}

// Mode returns file's permission mode.
func (s StatT) Mode() uint32 {
	return s.mode
}

// UID returns file's user id of owner.
func (s StatT) UID() uint32 {
	return s.uid
}

// GID returns file's group id of owner.
func (s StatT) GID() uint32 {
	return s.gid
}

// Rdev returns file's device ID (if it's special file).
func (s StatT) Rdev() uint64 {
	return s.rdev
}

// Size returns file's size.
func (s StatT) Size() int64 {
	return s.size
}

// Mtim returns file's last modification time.
func (s StatT) Mtim() syscall.Timespec {
	return s.mtim
}

// GetLastModification returns file's last modification time.
func (s StatT) GetLastModification() syscall.Timespec {
	return s.Mtim()
}
