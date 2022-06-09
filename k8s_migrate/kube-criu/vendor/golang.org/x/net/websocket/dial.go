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
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"crypto/tls"
	"net"
)

func dialWithDialer(dialer *net.Dialer, config *Config) (conn net.Conn, err error) {
	switch config.Location.Scheme {
	case "ws":
		conn, err = dialer.Dial("tcp", parseAuthority(config.Location))

	case "wss":
		conn, err = tls.DialWithDialer(dialer, "tcp", parseAuthority(config.Location), config.TlsConfig)

	default:
		err = ErrBadScheme
	}
	return
}
