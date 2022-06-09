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
package listeners

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/docker/go-connections/sockets"
)

// Init creates new listeners for the server.
func Init(proto, addr, socketGroup string, tlsConfig *tls.Config) (ls []net.Listener, err error) {
	switch proto {
	case "tcp":
		l, err := sockets.NewTCPSocket(addr, tlsConfig)
		if err != nil {
			return nil, err
		}
		ls = append(ls, l)
	case "unix":
		l, err := sockets.NewUnixSocket(addr, socketGroup)
		if err != nil {
			return nil, fmt.Errorf("can't create unix socket %s: %v", addr, err)
		}
		ls = append(ls, l)
	default:
		return nil, fmt.Errorf("Invalid protocol format: %q", proto)
	}

	return
}
