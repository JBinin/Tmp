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
// +build !linux,!windows,!freebsd,!solaris,!openbsd,!darwin

package system

import (
	"syscall"
)

// fromStatT creates a system.StatT type from a syscall.Stat_t type
func fromStatT(s *syscall.Stat_t) (*StatT, error) {
	return &StatT{size: s.Size,
		mode: uint32(s.Mode),
		uid:  s.Uid,
		gid:  s.Gid,
		rdev: uint64(s.Rdev),
		mtim: s.Mtimespec}, nil
}
