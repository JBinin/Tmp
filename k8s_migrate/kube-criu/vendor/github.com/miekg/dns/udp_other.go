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
// +build !linux,!plan9

package dns

import (
	"net"
	"syscall"
)

// These do nothing. See udp_linux.go for an example of how to implement this.

// We tried to adhire to some kind of naming scheme.

func setUDPSocketOptions4(conn *net.UDPConn) error                 { return nil }
func setUDPSocketOptions6(conn *net.UDPConn) error                 { return nil }
func getUDPSocketOptions6Only(conn *net.UDPConn) (bool, error)     { return false, nil }
func getUDPSocketName(conn *net.UDPConn) (syscall.Sockaddr, error) { return nil, nil }
