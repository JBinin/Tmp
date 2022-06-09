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

package xnet

import (
	"net"
	"time"
)

// ListenLocal opens a local socket for control communication
func ListenLocal(socket string) (net.Listener, error) {
	// on unix it's just a unix socket
	return net.Listen("unix", socket)
}

// DialTimeoutLocal is a DialTimeout function for local sockets
func DialTimeoutLocal(socket string, timeout time.Duration) (net.Conn, error) {
	// on unix, we dial a unix socket
	return net.DialTimeout("unix", socket, timeout)
}
