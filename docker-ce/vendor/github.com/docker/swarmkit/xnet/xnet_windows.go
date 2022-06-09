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

package xnet

import (
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

// ListenLocal opens a local socket for control communication
func ListenLocal(socket string) (net.Listener, error) {
	// set up ACL for the named pipe
	// allow Administrators and SYSTEM
	sddl := "D:P(A;;GA;;;BA)(A;;GA;;;SY)"
	c := winio.PipeConfig{
		SecurityDescriptor: sddl,
		MessageMode:        true,  // Use message mode so that CloseWrite() is supported
		InputBufferSize:    65536, // Use 64KB buffers to improve performance
		OutputBufferSize:   65536,
	}
	// on windows, our socket is actually a named pipe
	return winio.ListenPipe(socket, &c)
}

// DialTimeoutLocal is a DialTimeout function for local sockets
func DialTimeoutLocal(socket string, timeout time.Duration) (net.Conn, error) {
	// On windows, we dial a named pipe
	return winio.DialPipe(socket, &timeout)
}
