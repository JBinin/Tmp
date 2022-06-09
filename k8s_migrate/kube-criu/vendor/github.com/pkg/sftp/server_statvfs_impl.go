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
// +build darwin linux,!gccgo

// fill in statvfs structure with OS specific values
// Statfs_t is different per-kernel, and only exists on some unixes (not Solaris for instance)

package sftp

import (
	"syscall"
)

func (p sshFxpExtendedPacketStatVFS) respond(svr *Server) error {
	stat := &syscall.Statfs_t{}
	if err := syscall.Statfs(p.Path, stat); err != nil {
		return svr.sendPacket(statusFromError(p, err))
	}

	retPkt, err := statvfsFromStatfst(stat)
	if err != nil {
		return svr.sendPacket(statusFromError(p, err))
	}
	retPkt.ID = p.ID

	return svr.sendPacket(retPkt)
}
