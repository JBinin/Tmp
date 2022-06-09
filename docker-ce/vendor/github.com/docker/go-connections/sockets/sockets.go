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
// Package sockets provides helper functions to create and configure Unix or TCP sockets.
package sockets

import (
	"errors"
	"net"
	"net/http"
	"time"
)

// Why 32? See https://github.com/docker/docker/pull/8035.
const defaultTimeout = 32 * time.Second

// ErrProtocolNotAvailable is returned when a given transport protocol is not provided by the operating system.
var ErrProtocolNotAvailable = errors.New("protocol not available")

// ConfigureTransport configures the specified Transport according to the
// specified proto and addr.
// If the proto is unix (using a unix socket to communicate) or npipe the
// compression is disabled.
func ConfigureTransport(tr *http.Transport, proto, addr string) error {
	switch proto {
	case "unix":
		return configureUnixTransport(tr, proto, addr)
	case "npipe":
		return configureNpipeTransport(tr, proto, addr)
	default:
		tr.Proxy = http.ProxyFromEnvironment
		dialer, err := DialerFromEnvironment(&net.Dialer{
			Timeout: defaultTimeout,
		})
		if err != nil {
			return err
		}
		tr.Dial = dialer.Dial
	}
	return nil
}
