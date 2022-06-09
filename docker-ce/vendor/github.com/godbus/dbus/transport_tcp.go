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
//+build !windows

package dbus

import (
	"errors"
	"net"
)

func init() {
	transports["tcp"] = newTcpTransport
}

func tcpFamily(keys string) (string, error) {
	switch getKey(keys, "family") {
	case "":
		return "tcp", nil
	case "ipv4":
		return "tcp4", nil
	case "ipv6":
		return "tcp6", nil
	default:
		return "", errors.New("dbus: invalid tcp family (must be ipv4 or ipv6)")
	}
}

func newTcpTransport(keys string) (transport, error) {
	host := getKey(keys, "host")
	port := getKey(keys, "port")
	if host == "" || port == "" {
		return nil, errors.New("dbus: unsupported address (must set host and port)")
	}

	protocol, err := tcpFamily(keys)
	if err != nil {
		return nil, err
	}
	socket, err := net.Dial(protocol, net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}
	return NewConn(socket)
}
