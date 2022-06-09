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

package sockets

import (
	"net"
	"os"
	"syscall"
)

// NewUnixSocket creates a unix socket with the specified path and group.
func NewUnixSocket(path string, gid int) (net.Listener, error) {
	if err := syscall.Unlink(path); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	mask := syscall.Umask(0777)
	defer syscall.Umask(mask)

	l, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}
	if err := os.Chown(path, 0, gid); err != nil {
		l.Close()
		return nil, err
	}
	if err := os.Chmod(path, 0660); err != nil {
		l.Close()
		return nil, err
	}
	return l, nil
}
