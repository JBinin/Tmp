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
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"
)

const maxUnixSocketPathSize = len(syscall.RawSockaddrUnix{}.Path)

func configureUnixTransport(tr *http.Transport, proto, addr string) error {
	if len(addr) > maxUnixSocketPathSize {
		return fmt.Errorf("Unix socket path %q is too long", addr)
	}
	// No need for compression in local communications.
	tr.DisableCompression = true
	tr.Dial = func(_, _ string) (net.Conn, error) {
		return net.DialTimeout(proto, addr, defaultTimeout)
	}
	return nil
}

func configureNpipeTransport(tr *http.Transport, proto, addr string) error {
	return ErrProtocolNotAvailable
}

// DialPipe connects to a Windows named pipe.
// This is not supported on other OSes.
func DialPipe(_ string, _ time.Duration) (net.Conn, error) {
	return nil, syscall.EAFNOSUPPORT
}
